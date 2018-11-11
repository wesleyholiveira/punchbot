package configs

import "os"

// import "os"

// // Settings
// var Timer = os.Getenv("TIMER")

// // Discord
// var AuthURL = os.Getenv("AUTH_URL")

// var NotificationChannelsID = os.Getenv("NOTIFICATION_CHANNELS")
// var CommandChannelsID = os.Getenv("COMMAND_CHANNELS")
// var DiscordToken = os.Getenv("DISCORD_TOKEN")

// // Site
// var PunchEndpoint = os.Getenv("PUNCHSITE")
// var Home = PunchEndpoint + "/home"
// var Calendar = PunchEndpoint + "/calendario"

// // Redis
// var RedisHost = os.Getenv("REDIS_HOST")
// var RedisPort = os.Getenv("REDIS_PORT")
// var RedisPassword = os.Getenv("REDIS_PASSWORD")

var Timer = "15000"

// Discord
var AuthURL = "https://discordapp.com/api/oauth2/authorize?client_id=503005361419321345&permissions=0&scope=bot"

var NotificationChannelsID = "505535058233393153[@everyone]"
var CommandChannelsID = "455561651920568323"
var DiscordToken = "NTAzMDA1MzYxNDE5MzIxMzQ1.Dr1X3g.4_ldmPhe4BT-7hWiAfXp4AQr6LU"

// Site
var PunchEndpoint = "http://localhost:8080"
var Home = PunchEndpoint + "/home"
var Calendar = PunchEndpoint + "/calendario"

// Redis
var RedisHost = os.Getenv("REDIS_HOST")
var RedisPort = os.Getenv("REDIS_PORT")
var RedisPassword = os.Getenv("REDIS_PASSWORD")
