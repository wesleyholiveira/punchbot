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
													go notifyUser(s, p, myNots.Projects, key, userMention)
													break
												} else {
													s.ChannelMessageSend(key,
														fmt.Sprintf("Vejo que você possui interesse em receber notificações "+
															"via DM mas isto é **exclusivo** para usuários **VIP** no Discord.\n"+
															"Leia o canal #regras ou acesse: %s para adquirir seu **VIP!**", configs.PunchEndpoint))
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

func notify(s *discordgo.Session, current *[]models.Project, prev *[]models.Project, channelID, userMention string) (bool, error) {
	punchReleases := models.GetProjects()
	cLen := len(*current)
	pLen := len(*prev)
	diff := int(math.Abs(float64(cLen) - float64(pLen)))

	if diff == 0 {
		diff = 1
	}

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d", diff, pLen, cLen)

	currentSlice := (*current)[0:diff]

	for i, c := range currentSlice {
		if !c.AlreadyReleased {
			for _, p := range *prev {
				if c.IDProject != p.IDProject {
					sendMessage(s, &c, channelID, userMention)
					(*current)[i].AlreadyReleased = true
					break
				}
			}
		} else {
			log.Warnf("A previosly released project %s[%s] was already released.",
				c.Project, c.HashID)
			log.Warn("Ignoring...")
			break
		}
	}

	*punchReleases = *current

	return true, nil
}

func notifyUser(s *discordgo.Session, current *[]models.Project, prev *[]models.Project, channelID, userMention string) (bool, error) {
	punchReleases := models.GetProjects()
	cLen := len(*current)
	pLen := len(*prev)
	diff := int(math.Abs(float64(cLen) - float64(pLen)))

	if diff == 0 {
		diff = 1
	}

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d", diff, pLen, cLen)
	currentSlice := (*current)[0:diff]

	for i, c := range currentSlice {
		if !c.AlreadyReleased {
			for _, p := range *prev {
				if c.IDProject == p.IDProject {
					sendMessage(s, &c, channelID, userMention)
					(*current)[i].AlreadyReleased = true
					break
				}
			}
		} else {
			log.Warnf("A previosly released project %s[%s] was already released.",
				c.Project, c.HashID)
			log.Warn("Ignoring...")
			break
		}
	}

	*punchReleases = *current

	return true, nil
}

func sendMessage(s *discordgo.Session, c *models.Project, channelID, userMention string) (bool, error) {
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
		log.Errorf("[%s] Image error: %s", c.Project, err)
	}

	return true, nil
}
