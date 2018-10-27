package parallelism

import (
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
)

func Notifier(s *discordgo.Session, projects chan *[]models.Project) {
	for p := range projects {
		log.Info("Entrou no notifier")

		var msg string
		myNotifications := models.GetNotifyUser()

		if myNotifications != nil {
			for key := range myNotifications {
				prjs := *myNotifications[key].Projects
				myProjects := make([]models.Project, len(*p))
				diff := int(math.Abs(float64(len(*p) - len(myProjects))))
				copy(myProjects, prjs)

				log.Infof("Diff: %d, PUNCH PROJECTS: %d, MY PROJECTS: %d", diff, len(*p), len(myProjects))

				for _, project := range *p {
					for i := 0; i < diff; i++ {
						if project.IDProject == myProjects[i].IDProject {
							msg += fmt.Sprintf("O **anime %s** acabou de sair no site da **__%s__**", project.IDProject, configs.PunchEndpoint)
						}
					}
				}
			}
		}
	}
}
