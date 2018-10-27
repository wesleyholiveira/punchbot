package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/wesleyholiveira/punchbot/models"
)

var cmds models.Commands
var ProjectChan chan *[]models.Project

func init() {
	cmds = make(models.Commands)
	ProjectChan = make(chan *[]models.Project, 0)

	cmds["list"] = List
	cmds["notify"] = Notify
}

func Entry(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := "p!"
	if strings.HasPrefix(m.Content, prefix) {
		projects := models.GetProjects()
		m.Content = m.Content[2:]
		command := strings.Split(m.Content, " ")
		cmd := command[0]
		if _, ok := cmds[cmd]; ok {
			args := command[1:]
			cmds[cmd](*projects, s, m, args)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Comando não encontrado.")
		}
	}
}
