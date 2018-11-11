package commands

import (
	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	msgHelp := `
**p!list<,...animes>** - Retorna uma lista de IDs compativel com o nome digitado. __(ex: p!list Boruto)__
**p!notify<,...ids>** - Notifica via DM os animes selecionados. __(ex: p!notify 2688 8)__
**p!auth** - Retorna a URL de autorização para interagir com o bot via DM.
**p!help** - Retorna os comandos disponíveis`
	s.ChannelMessageSend(m.ChannelID, msgHelp)
}
