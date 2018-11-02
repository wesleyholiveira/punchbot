package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/wesleyholiveira/punchbot/configs"
)

func Auth(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Para receber notificações via DM, acesse o link: %s\n", configs.AuthURL))
}
