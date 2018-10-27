package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
)

var client *redis.Client

func NewClient() *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.RedisHost, configs.RedisPort),
		Password: configs.RedisPassword, // no password set
		DB:       0,                     // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Errorln("Redis", err)
	}

	log.Debug(pong)
	return client
}

func GetClient() *redis.Client {
	return client
}
