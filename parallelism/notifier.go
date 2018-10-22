package parallelism

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/wesleyholiveira/punchbot/models"
)

func Notifier(s *discordgo.Session, notify chan bool) {
	for n := range notify {
		fmt.Println("Entrou no notifier", n)
		for key, value := range models.GetNotifyUser() {
			fmt.Printf("%s %#v\n", key, value)
		}
	}
}
