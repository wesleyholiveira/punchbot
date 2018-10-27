package commands

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/redis"

	"github.com/bwmarrin/discordgo"
	"github.com/wesleyholiveira/punchbot/models"
)

func Notify(projects []models.Project, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	redis := redis.GetClient()
	notifyUser := models.GetNotifyUser()
	notifyRedis := &models.Notify{}

	redisGet, err := redis.Get(m.ChannelID).Result()
	if err != nil {
		log.Errorln("Redis GET", err)
	}

	if len(args) == 0 {

		if len(redisGet) > 0 {
			err = json.Unmarshal([]byte(redisGet), notifyRedis)

			if err != nil {
				log.Errorln("Unmarshal", err)
			}

			notifyUser[m.ChannelID] = notifyRedis
		}

		msg := "Sua lista de animes a serem notificados:\n"

		for key, _ := range notifyUser {
			for _, prj := range *notifyUser[key].Projects {
				msg += fmt.Sprintf("**%s** - %s\n", prj.IDProject, prj.Project)
			}
		}

		s.ChannelMessageSend(m.ChannelID, msg)
	} else {
		projectsUser := make([]models.Project, 0, len(args))

		for _, project := range projects {
			for _, arg := range args {
				if project.IDProject == arg {
					projectsUser = append(projectsUser, project)
				}
			}
		}

		notify := &models.Notify{UserID: m.ChannelID, Projects: &projectsUser}

		notifyUser[notify.UserID] = notify
		ProjectChan <- &projectsUser

		notifyMarshal, err := json.Marshal(notify)
		if err != nil {
			log.Errorln("Marshal", err)
		}

		notifyMarshalStr := string(notifyMarshal)
		err = redis.Set(notify.UserID, notifyMarshalStr, 0).Err()

		if err != nil {
			log.Errorln("Redis SET", err)
		}

		err = redis.Save().Err()
		if err != nil {
			log.Errorln("Redis SAVE", err)
		}

		s.ChannelMessageSend(m.ChannelID, "Você será notificado quando um episódio novo dos itens selecionados aparecerem no site.")
	}
}
