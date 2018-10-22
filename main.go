package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/commands"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/parallelism"
)

func main() {
	timer := time.NewTicker(15 * time.Second)
	notifyChan := make(chan bool, 1)

	log.SetOutput(os.Stdout)

	d, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Error(err)
	}

	err = d.Open()
	if err != nil {
		log.Error(err)
	}

	d.AddHandler(commands.Entry)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go parallelism.TickerHTTP(timer, notifyChan)
	go parallelism.Notifier(d, notifyChan)

	<-sc
	defer d.Close()
}
