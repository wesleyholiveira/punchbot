package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/wesleyholiveira/punchbot/configs"
	"golang.org/x/oauth2"
)

var client *twitter.Client

func NewClient() (*twitter.Client, error) {
	config := oauth1.NewConfig(configs.TwitterConsumer, configs.TwitterConsumerSecret)
	token := oauth1.NewToken(configs.TwitterAccessToken, configs.TwitterAccessTokenSecret)
	httpClient := config.Client(oauth2.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)

	client = twitterClient
	return twitterClient, nil
}

func GetClient() *twitter.Client {
	return client
}
