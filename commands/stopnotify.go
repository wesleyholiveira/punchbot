package commands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/redis"
)

func StopNotify(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	channel := m.ChannelID
	redis := redis.GetClient()
	notifyUser := models.GetNotifyUser()
	notifyRedis := &models.Notify{}

	_, err := redis.Del(channel).Result()
	if err != nil {
		log.Error("Redis: ", err)
	}

	err = redis.Save().Err()
	if err != nil {
		log.Error("Redis: ", err)
	}

	notifyUser[channel] = notifyRedis
	s.ChannelMessageSend(channel, "Você não receberá mais notificações via DM.")
}
