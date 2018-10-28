package configs

import (
	"os"
	"time"
)

// Settings
var Timer time.Duration = 15000

// Discord
var AuthURL = os.Getenv("AUTH_URL")

var ChannelID = os.Getenv("DISCORD_CHANNEL")
var DiscordToken = os.Getenv("DISCORD_TOKEN")

// Site
var PunchEndpoint = os.Getenv("PUNCHSITE")
var Home = PunchEndpoint + "/home"
var Calendar = PunchEndpoint + "/calendario"

// Redis
var RedisHost = os.Getenv("REDIS_HOST")
var RedisPort = os.Getenv("REDIS_PORT")
var RedisPassword = os.Getenv("REDIS_PASSWORD")
