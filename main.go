package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wesleyholiveira/punchbot/redis"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/commands"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/parallelism"
)

func main() {
	rClient := redis.NewClient()
	timer := time.NewTicker(15 * time.Second)

	log.SetOutput(os.Stdout)

	d, err := discordgo.New("Bot " + configs.DiscordToken)
	if err != nil {
		log.Errorln("Discord", err)
	}

	err = d.Open()
	if err != nil {
		log.Errorln("Discord", err)
	}

	d.AddHandler(commands.Entry)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go parallelism.TickerHTTP(timer, commands.ProjectChan)
	go parallelism.Notifier(d, commands.ProjectChan)

	<-sc
	defer rClient.Close()
	defer d.Close()
}
