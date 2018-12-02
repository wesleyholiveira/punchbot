package commands

import (
	"encoding/json"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/redis"
)

func Hyped(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var msg string

	var allowedUsers []string
	allowedUsersArray := strings.Split(",", configs.AllowedUsers)

	if len(allowedUsersArray) < 2 && allowedUsersArray[0] == "," {
		allowedUsers = append(allowedUsers, configs.AllowedUsers)
	} else {
		allowedUsers = allowedUsersArray
	}

	msg = "Você não tem permissão para utilizar este comando."

	for _, userID := range allowedUsers {
		if m.Author.ID == userID {
			r := redis.GetClient()
			calendar := *models.GetCalendarProjects()

			if len(args) == 0 {
				fbWhiteList, err := r.Get("fbWhiteList").Result()

				if err != nil {
					msg = "Nenhum anime HYPED foi salvo."
				} else {
					fbProjects := make([]models.Project, 0, cap(calendar))
					err = json.Unmarshal([]byte(fbWhiteList), &fbProjects)

					if err != nil {
						log.Error(err)
						return
					}

					msg = AnimeList(&fbProjects, func(project *models.Project) bool { return true })
				}
			} else {
				msg = "Os animes HYPEDS foram salvos com sucesso!"
				fbProjects := make([]models.Project, 0, cap(args))

				for _, arg := range args {
					for _, item := range calendar {
						projectID := arg

						if projectID == item.IDProject {
							fbProjects = append(fbProjects, item)
						}
					}
				}

				data, err := json.Marshal(fbProjects)
				if err != nil {
					log.Error(err)
					return
				}

				_, err = r.Set("fbWhiteList", string(data), 0).Result()
				if err != nil {
					log.Error(err)
					return
				}

				r.Save()
			}
			break
		}
	}

	s.ChannelMessageSend(m.ChannelID, msg)
}
