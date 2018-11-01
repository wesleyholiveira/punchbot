package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/wesleyholiveira/punchbot/configs"
)

func Auth(s *discordgo.Session, channel string, args []string) {
	s.ChannelMessageSend(channel, fmt.Sprintf("Para receber notificações via DM, acesse o link: %s\n", configs.AuthURL))
}
