package parallelism

import (
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/helpers"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

var channels map[string]string

func init() {
	channels = helpers.ParseChannels(configs.NotificationChannelsID)
}

func Notifier(s *discordgo.Session, projects chan *[]models.Project) {
	for p := range projects {
		var guildID string
		prev := models.GetProjects()
		log.Info("Notifier is on")

		for key, tag := range channels {

			ch, _ := s.Channel(key)

			if ch != nil {
				guild, _ := s.Guild(ch.GuildID)
				guildID = guild.ID

				userMention := ""
				if tag != "" {
					username := tag[1:]
					if username != "everyone" {
						for _, m := range guild.Members {
							if m.User.Username == username {
								userMention = m.User.Mention()
								break
							}
						}
						if userMention == "" {
							for _, r := range guild.Roles {
								if r.Name == username && r.Mentionable {
									userMention = fmt.Sprintf("<@&%s>", r.ID)
									break
								}
							}
						}
					} else {
						userMention = "@" + username
					}
					userMention += " "
				}
				go notify(s, p, prev, ch.ID, userMention)
			}
		}

		myNotifications := models.GetNotifyUser()
		if myNotifications != nil {
			guild, _ := s.Guild(guildID)
			userMention := ""
			for key := range myNotifications {
				ch, _ := s.Channel(key)
				myNots := myNotifications[key]

				if ch != nil {
					if guild != nil {
						if ch.Type == discordgo.ChannelTypeDM {
							if myNots != nil {
								for _, m := range guild.Members {
									if myNots.UserID == m.User.ID {
										log.Info("User found in punch's server")
										for _, userRoleID := range m.Roles {
											for _, role := range guild.Roles {
												if role.Name == "VIP" && role.ID == userRoleID {
													log.Info("The user is a vip!!")
													log.Info("Sending notifications (if exists)")
													go notify(s, myNots.Projects, prev, key, userMention)
													break
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
		}
	}
}

func notify(s *discordgo.Session, current *[]models.Project, prev *[]models.Project, channelID, userMention string) error {
	cLen := len(*current)
	pLen := len(*prev)
	diff := int(math.Abs(float64(cLen) - float64(pLen)))

	if diff == 0 {
		diff = 1
	}

	currentSlice := (*current)[0:diff]
	prevSlice := (*prev)[0:diff]

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d", diff, pLen, cLen)

	for _, c := range currentSlice {
		for _, p := range prevSlice {
			if c.IDProject != p.IDProject {
				log.Info("PROJECT MATCHED!")

				screen := c.Screen
				img := strings.Split(screen, "/")
				imgName := img[len(img)-1]

				if !strings.Contains(c.Screen, "http") {
					screen = configs.PunchEndpoint + c.Screen
				}

				httpImage, err := services.Get(screen)

				if err == nil {

					respImage := httpImage.Body
					defer respImage.Close()

					msg := fmt.Sprintf("%sO **%s** do anime **%s** acabou de ser lançado! -> %s\n",
						userMention,
						c.Number,
						c.Project,
						configs.PunchEndpoint+c.Link)

					_, err := s.ChannelFileSendWithMessage(channelID, msg, imgName, respImage)
					if err != nil {
						log.Error(err)
					}
				} else {
					log.Error("Image error: ", err)
				}
			}
		}
	}
	return nil
}
