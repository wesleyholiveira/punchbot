package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/redis"
)

func Notify(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	channel := m.ChannelID
	redis := redis.GetClient()
	notifyUser := models.GetNotifyUser()
	notifyRedis := &models.Notify{}

	redisGet, err := redis.Get(channel).Result()
	if err != nil {
		log.Errorln("Redis GET", err)
	}

	if len(args) == 0 || args[0] == "" {

		if len(redisGet) > 0 {
			err = json.Unmarshal([]byte(redisGet), notifyRedis)

			if err != nil {
				log.Errorln("Unmarshal", err)
			}

			notifyUser[channel] = notifyRedis
		}

		msg := "Não há nada para ser listado."

		for key := range notifyUser {
			if nUser := notifyUser[key].Projects; nUser != nil {
				msg = "Sua lista de animes a serem notificados:\n"
				for _, prj := range *nUser {
					msg += fmt.Sprintf("**%s** - %s\n", prj.IDProject, prj.Project)
				}
			}
		}

		s.ChannelMessageSend(channel, msg)
	} else {
		projects := models.GetCalendarProjects()
		projectsUser := make([]models.Project, 0, len(args))

		for _, project := range *projects {
			for _, arg := range args {
				p := strings.TrimSpace(strings.ToLower(project.Project))
				c := strings.TrimSpace(strings.ToLower(arg))

				if strings.Contains(p, c) {
					projectsUser = append(projectsUser, project)
				}
			}
		}

		notify := &models.Notify{UserID: m.Author.ID, Projects: &projectsUser}

		notifyUser[channel] = notify
		notifyMarshal, err := json.Marshal(notify)

		if err != nil {
			log.Errorln("Marshal", err)
		}

		notifyMarshalStr := string(notifyMarshal)
		err = redis.Set(channel, notifyMarshalStr, 0).Err()

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
