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
var DiscordToken = "NTAzMDA1MzYxNDE5MzIxMzQ1.DrvaHA.TL8_hp9cWAIVVPz66VIf9JlbXvc"

// Site
var PunchEndpoint = "https://punchsubs.net"
var Home = PunchEndpoint + "/home"
var Calendar = PunchEndpoint + "/calendario"

// Redis
var RedisHost = os.Getenv("REDIS_HOST")
var RedisPort = os.Getenv("REDIS_PORT")
var RedisPassword = os.Getenv("REDIS_PASSWORD")
