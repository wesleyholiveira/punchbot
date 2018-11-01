package commands

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/redis"
)

func Notify(s *discordgo.Session, channel string, args []string) {
	redis := redis.GetClient()
	notifyUser := models.GetNotifyUser()
	notifyRedis := &models.Notify{}

	redisGet, err := redis.Get(channel).Result()
	if err != nil {
		log.Errorln("Redis GET", err)
	}

	if len(args) == 0 {

		if len(redisGet) > 0 {
			err = json.Unmarshal([]byte(redisGet), notifyRedis)

			if err != nil {
				log.Errorln("Unmarshal", err)
			}

			notifyUser[channel] = notifyRedis
		}

		msg := "Sua lista de animes a serem notificados:\n"

		for key := range notifyUser {
			for _, prj := range *notifyUser[key].Projects {
				msg += fmt.Sprintf("**%s** - %s\n", prj.IDProject, prj.Project)
			}
		}

		s.ChannelMessageSend(channel, msg)
	} else {
		projects := models.GetCalendarProjects()
		projectsUser := make([]models.Project, 0, len(args))

		for _, project := range *projects {
			for _, arg := range args {
				if project.IDProject == arg {
					projectsUser = append(projectsUser, project)
				}
			}
		}

		notify := &models.Notify{UserID: channel, Projects: &projectsUser}

		notifyUser[notify.UserID] = notify
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

		s.ChannelMessageSend(channel, "Você será notificado quando um episódio novo dos itens selecionados aparecerem no site.")
	}
}
