package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/parallelism"
)

var sinCache map[string]models.Project

func init() {
	sinCache = make(map[string]models.Project)
}

func DescriptionList(projects *[]models.Project, callback sCallback) *discordgo.MessageEmbed {
	var extraInfosProject models.Project
	for _, project := range *projects {
		if callback(&project) {
			if c, ok := sinCache[project.IDProject]; !ok {
				project.AlreadyReleased = false
				prj := make([]models.Project, 1)
				prj[0] = project
				extraInfosProject = parallelism.GetExtraInfos(&prj)[0]
				sinCache[project.IDProject] = extraInfosProject
			} else {
				extraInfosProject = c
			}

			for _, p := range *projects {
				if p.ID == project.ID {
					project.Description = extraInfosProject.Description
					project.ExtraInfos = extraInfosProject.ExtraInfos
				}
			}
			return descriptionFormat(project.Project, project.Description, project.Link)
		}
	}
	return nil
}

func descriptionFormat(projectName, desc, url string) *discordgo.MessageEmbed {
	link := configs.PunchEndpoint + url
	icon := fmt.Sprintf("%s/imagens/favicon-96x96.png", configs.PunchEndpoint)

	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "PUNCH! Fansubs",
			IconURL: icon,
		},
		Title:       fmt.Sprintf("%s", projectName),
		Description: fmt.Sprintf("%s", desc),
		URL:         link,
		Color:       65280,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "PUNCH! Fansubs",
			IconURL: icon,
		},
	}
}

func Sinopse(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	projects := models.GetProjects()
	list := &discordgo.MessageEmbed{}

	if len(args) > 0 && args[0] != "" {
		list = DescriptionList(projects, func(project *models.Project) bool {
			for _, arg := range args {
				lArg := strings.ToLower(arg)
				contains := strings.Contains(strings.ToLower(project.Project), lArg)
				if contains {
					return true
				}
			}
			return false
		})
	}

	if list == nil {
		s.ChannelMessageSend(m.ChannelID, "Projeto não encontrado.")
		return
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, list)
	if err != nil {
		log.Error(err)
	}
}
