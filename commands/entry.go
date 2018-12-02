package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
)

var cmds models.Commands
var ProjectChan, PrevProjectChan chan *[]models.Project
var except []string

func init() {
	cmds = make(models.Commands)
	ProjectChan = make(chan *[]models.Project, 0)
	PrevProjectChan = make(chan *[]models.Project, 0)

	except = []string{"help", "auth", "list"}

	cmds["list"] = List
	cmds["hyped"] = Hyped
	cmds["notify"] = Notify
	cmds["stopnotify"] = StopNotify
	cmds["auth"] = Auth
	cmds["help"] = Help
}

func Entry(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := "p!"
	if strings.HasPrefix(m.Content, prefix) {
		m.Content = m.Content[2:]
		command := strings.Split(m.Content, " ")
		cmd := command[0]
		if _, ok := cmds[cmd]; ok {
			ch, _ := s.Channel(m.ChannelID)
			guild, _ := s.Guild(ch.GuildID)
			args := command[1:]

			if guild == nil {
				cmds[cmd](s, m, args)
			} else if isException(cmd) &&
				strings.Contains(m.ChannelID, configs.CommandChannelsID) {
				cmds[cmd](s, m, args)
			}

		} else {
			s.ChannelMessageSend(m.ChannelID, "Comando não encontrado.")
		}
	}
}

func isException(exception string) bool {
	for _, val := range except {
		if val == exception {
			return true
		}
	}
	return false
}
