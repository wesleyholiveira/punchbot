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
	"github.com/wesleyholiveira/punchbot/redis"
)

func main() {
	rClient := redis.NewClient()
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	now := time.Now().In(loc)
	timer := time.NewTicker(configs.Timer * time.Millisecond)

	log.SetOutput(os.Stdout)
	log.WithTime(now)

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

	go parallelism.GetProjects()
	go parallelism.GetProjectsCalendar()
	go parallelism.TickerHTTP(timer, commands.ProjectChan)
	go parallelism.Notifier(d, commands.ProjectChan)

	<-sc
	defer rClient.Close()
	defer d.Close()
}
