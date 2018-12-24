package parallelism

import (
	"encoding/json"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/models"
)

func Notify(redis *redis.Client) {
	log.Info("Getting the personal notifications...")

	notifyUser := models.GetNotifyUser()
	notifyRedis := &models.Notify{}
	redisKeys, err := redis.Keys("*").Result()

	if err != nil {
		log.Error(err)
	}

	for _, channel := range redisKeys {
		redisGet, err := redis.Get(channel).Result()
		if err != nil {
			log.Errorln("Redis GET", err)
		}

		err = json.Unmarshal([]byte(redisGet), notifyRedis)

		if err != nil {
			log.Errorln("Unmarshal", err)
		}

		notifyUser[channel] = notifyRedis
	}
}
