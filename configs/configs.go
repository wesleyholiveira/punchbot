package configs

import "os"

// Settings
var ProxyURL = os.Getenv("PROXY_URL")
var AllowedUsers = os.Getenv("PUNCHBOT_ALLOWERDUSERS_HYPED")
var Timer = os.Getenv("TIMER")

// Discord
var AuthURL = os.Getenv("AUTH_URL")

var NotificationChannelsID = os.Getenv("NOTIFICATION_CHANNELS")
var CommandChannelsID = os.Getenv("COMMAND_CHANNELS")
var DiscordToken = os.Getenv("DISCORD_TOKEN")

// Twitter
var TwitterConsumer = os.Getenv("TWITTER_CONSUMER")
var TwitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
var TwitterAccessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
var TwitterAccessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

// Facebook
var FacebookAppID = os.Getenv("FACEBOOK_APPID")
var FacebookAppSecret = os.Getenv("FACEBOOK_APPID_SECRET")
var FacebookCallbackURL = "https://localhost/"

// Site
var PunchEndpoint = os.Getenv("PUNCHSITE")
var Home = PunchEndpoint + "/home"
var Calendar = PunchEndpoint + "/calendario"

// Redis
var RedisHost = os.Getenv("REDIS_HOST")
var RedisPort = os.Getenv("REDIS_PORT")
var RedisPassword = os.Getenv("REDIS_PASSWORD")
