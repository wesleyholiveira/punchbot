package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/models"
)

type sCallback func(prj *models.Project) bool

func AnimeList(projects *[]models.Project, callback sCallback) string {
	list := ""
	for _, project := range *projects {
		if callback(&project) {
			list += animeListFormat(project.IDProject, project.Project, project.Day)
		}
	}
	return list
}

func animeListFormat(id, projectName, day string) string {
	return fmt.Sprintf("**%s** - %s (**%s**)\n", id, projectName, day)
}

func List(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	projects := models.GetCalendarProjects()
	list := ""

	if len(args) > 0 && args[0] != "" {
		list = AnimeList(projects, func(project *models.Project) bool {
			for _, arg := range args {
				if strings.Contains(strings.ToLower(project.Project), strings.ToLower(arg)) {
					return true
				}
			}
			return false
		})
	} else if len(*projects) > 0 {
		list = AnimeList(projects, func(project *models.Project) bool { return true })
		rows := strings.Split(list, "\n")
		list = strings.Join(rows[:9], "\n")
	}

	if len(list) < 1 {
		list = "Projeto não encontrado."
	}

	_, err := s.ChannelMessageSend(m.ChannelID, list)
	if err != nil {
		log.Error(err)
	}
}
