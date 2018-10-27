package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/models"
)

type sCallback func(prj *models.Project) bool

func animeList(projects *[]models.Project, callback sCallback) string {
	list := ""
	for _, project := range *projects {
		if callback(&project) {
			list += animeListFormat(project.IDProject, project.Project)
		}
	}
	return list
}

func animeListFormat(id, projectName string) string {
	return fmt.Sprintf("**%s** - %s\n", id, projectName)
}

func List(projects []models.Project, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	list := ""

	if len(args) > 0 {
		list = animeList(&projects, func(project *models.Project) bool {
			for _, arg := range args {
				if strings.Contains(strings.ToLower(project.Project), strings.ToLower(arg)) {
					return true
				}
			}
			return false
		})
	} else {
		list = animeList(&projects, func(project *models.Project) bool { return true })
		rows := strings.Split(list, "\n")
		list = strings.Join(rows[0:9], "\n")
	}

	_, err := s.ChannelMessageSend(m.ChannelID, list)
	if err != nil {
		log.Error(err)
	}

}
