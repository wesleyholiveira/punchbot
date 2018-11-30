package facebook

import (
	fb "github.com/huandu/facebook"
	"github.com/wesleyholiveira/punchbot/configs"
	"golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook"
)

func NewClient() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     configs.FacebookAppID,
		ClientSecret: configs.FacebookAppSecret,
		RedirectURL:  configs.FacebookCallbackURL,
		Scopes:       []string{"email", "manage_pages", "publish_pages"},
		Endpoint:     oauth2fb.Endpoint,
	}

	return conf
}

func GetClient(conf *oauth2.Config, code string) (*fb.Session, error) {
	token, err := conf.Exchange(oauth2.NoContext, code)

	if err != nil {
		return nil, err
	}

	// Create a client to manage access token life cycle.
	client := conf.Client(oauth2.NoContext, token)

	// Use OAuth2 client with session.
	session := &fb.Session{
		Version:    "v3.1",
		HttpClient: client,
	}
	return session, nil
}

func GetClientByToken(token string) *fb.Session {
	globalApp := fb.New(configs.FacebookAppID, configs.FacebookAppSecret)
	globalApp.RedirectUri = configs.FacebookCallbackURL

	session := globalApp.Session(token)
	return session
}
