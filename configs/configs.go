package configs

import "os"

// Settings
var Timer = os.Getenv("TIMER")

// Discord
var AuthURL = os.Getenv("AUTH_URL")

var NotificationChannelsID = os.Getenv("NOTIFICATION_CHANNELS")
var CommandChannelsID = os.Getenv("COMMAND_CHANNELS")
var DiscordToken = os.Getenv("DISCORD_TOKEN")
var TwitterConsumer = os.Getenv("TWITTER_CONSUMER")
var TwitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
var TwitterAccessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
var TwitterAccessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

// Site
var PunchEndpoint = os.Getenv("PUNCHSITE")
var Home = PunchEndpoint + "/home"
var Calendar = PunchEndpoint + "/calendario"

// Redis
var RedisHost = os.Getenv("REDIS_HOST")
var RedisPort = os.Getenv("REDIS_PORT")
var RedisPassword = os.Getenv("REDIS_PASSWORD")
