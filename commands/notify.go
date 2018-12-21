package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/parallelism"
	"github.com/wesleyholiveira/punchbot/redis"
)

func Notify(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	c, _ := s.Channel(m.ChannelID)
	t := c.Type
	vip := false
	channel := m.ChannelID
	redis := redis.GetClient()
	notifyUser := models.GetNotifyUser()
	notifyRedis := &models.Notify{}

	user := m.Author
	if _, ok := notifyUser[channel]; !ok {
		for key, _ := range parallelism.Channels {
			ch, _ := s.Channel(key)
			guild, _ := s.Guild(ch.GuildID)

			if ch != nil {
				if guild != nil {
					if t == discordgo.ChannelTypeDM {
						for _, m := range guild.Members {
							if m.User.ID == user.ID {
								log.Infof("User found %s in punch's server", user.Username)
								for _, userRoleID := range m.Roles {
									for _, role := range guild.Roles {
										if role.Name == "VIP" && role.ID == userRoleID {
											log.Infof("The user %s is a vip!!", m.User.Username)
											vip = true
										} else {
											vip = false
										}
									}
								}
								break
							}
						}
					}
				}
			}
		}
	}

	if vip {
		if len(args) == 0 || args[0] == "" {

			redisGet, err := redis.Get(channel).Result()
			if err != nil {
				log.Errorln("Redis GET", err)
			}

			if len(redisGet) > 0 {
				err = json.Unmarshal([]byte(redisGet), notifyRedis)

				if err != nil {
					log.Errorln("Unmarshal", err)
				}

				notifyUser[channel] = notifyRedis
			}

			msg := "Não há nada para ser listado."

			if user := notifyUser[channel]; user != nil {
				nUser := user.Projects
				msg = "Sua lista de animes a serem notificados:\n"
				for _, prj := range *nUser {
					msg += fmt.Sprintf("**%s** - %s\n", prj.IDProject, prj.Project)
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

			notify := &models.Notify{UserID: user.ID, Projects: &projectsUser}
			notify.VIP = true
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
	} else {
		log.Infof("The user %s isn't a vip!!", user.Username)
		s.ChannelMessageSend(channel,
			fmt.Sprintf("Vejo que você possui interesse em receber notificações "+
				"via DM mas isto é **exclusivo** para usuários **VIP** no Discord.\n"+
				"Leia o canal #regras ou acesse: %s para adquirir seu **VIP!**",
				configs.PunchEndpoint))
	}

}
