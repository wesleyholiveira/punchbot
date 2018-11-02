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

var userMention string
var channels map[string]string

func init() {
	channels = helpers.ParseChannels(configs.NotificationChannelsID)
}

func Notifier(s *discordgo.Session, projects chan *[]models.Project) {
	for p := range projects {
		log.Info("Notifier is on")

		for key, tag := range channels {

			ch, _ := s.Channel(key)

			if ch != nil {
				guild, _ := s.Guild(ch.GuildID)

				userMention = ""
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

				if guild != nil {
					projectsPunch := models.GetProjects()
					notify(s, p, PrevProject, *projectsPunch, ch.ID)
				}
			}
		}

		myNotifications := models.GetNotifyUser()
		if myNotifications != nil {
			userMention = ""
			for key := range myNotifications {
				prjs := *myNotifications[key].Projects
				notify(s, p, PrevProject, prjs, key)
			}
		}
	}
}

func notify(s *discordgo.Session, p *[]models.Project, prev *[]models.Project, prjs []models.Project, channelID string) error {
	myProjects := make([]models.Project, len(prjs))
	diff := int(math.Abs(float64(len(*p) - len(*prev))))

	if diff == 0 {
		diff = 1
	}

	copy(myProjects, prjs)

	projectsSlice := *p
	projectsSlice = projectsSlice[0:diff]

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d, MY PROJECTS: %d", diff, len(*prev), len(*p), len(myProjects))

	for _, myProject := range myProjects {
		for _, project := range projectsSlice {
			if project.IDProject == myProject.IDProject {
				log.Info("PROJECT MATCHED!")
				screen := project.Screen
				img := strings.Split(screen, "/")
				imgName := img[len(img)-1]

				if !strings.Contains(project.Screen, "http") {
					screen = configs.PunchEndpoint + project.Screen
				}

				httpImage, err := services.Get(screen)

				if err == nil {

					respImage := httpImage.Body
					defer respImage.Close()

					msg := fmt.Sprintf("%sO **%s** do anime **%s** acabou de ser lançado! -> %s\n",
						userMention,
						project.Numero,
						project.Project,
						configs.PunchEndpoint+project.Link)

					_, err := s.ChannelFileSendWithMessage(channelID, msg, imgName, respImage)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
	return nil
}
