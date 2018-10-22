package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	model "github.com/wesleyholiveira/punchbot/models"
)

func List(projects []model.Project, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	animeList := ""

	for _, project := range projects {
		animeList += fmt.Sprintf("**%s** - %s\n", project.IDProject, project.Project)
	}

	s.ChannelMessageSend(m.ChannelID, animeList)
}
