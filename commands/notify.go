package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/wesleyholiveira/punchbot/models"
)

func Notify(projects []models.Project, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	notifyUser := models.GetNotifyUser()
	projectsUser := make([]models.Project, 0, len(args))
	notify := &models.Notify{UserID: m.Author.ID, Projects: &projectsUser}

	for _, project := range projects {
		for _, arg := range args {
			if project.IDProject == arg {
				projectsUser = append(projectsUser, project)
			}
		}
	}
	notifyUser[notify.UserID] = notify
}
