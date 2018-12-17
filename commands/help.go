package commands

import (
	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	msgHelp := `
**p!list<,...nome>** - Retorna os animes do calendário. Digite o nome para filtrar. __(ex: p!list Boruto, p!list Boruto, Slime)__
**PRIVADO** **p!notify<,...nome>** - Notifica via DM os animes selecionados pelo nome. __(ex: p!notify Boruto, Slime)__
**PRIVADO** **p!stopnotify** - Revoga o interesse em receber notificações via DM.
**p!auth** - Retorna a URL de autorização para interagir com o bot via DM.
**p!help** - Retorna os comandos disponíveis`
	s.ChannelMessageSend(m.ChannelID, msgHelp)
}
