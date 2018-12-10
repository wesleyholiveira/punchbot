package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/wesleyholiveira/punchbot/models"
	"golang.org/x/oauth2"

	"github.com/bwmarrin/discordgo"
	fb "github.com/huandu/facebook"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/commands"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/parallelism"
	"github.com/wesleyholiveira/punchbot/redis"
	"github.com/wesleyholiveira/punchbot/services/facebook"
	"github.com/wesleyholiveira/punchbot/services/twitter"
)

var fbOauth *oauth2.Config

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	rClient := redis.NewClient()
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	now := time.Now().In(loc)
	timeDuration, err := strconv.Atoi(configs.Timer)

	if err != nil {
		log.Errorln(err)
	}

	timer := time.NewTicker(time.Duration(timeDuration) * time.Millisecond)

	log.SetOutput(os.Stdout)
	log.WithTime(now)

	d, err := discordgo.New("Bot " + configs.DiscordToken)
	if err != nil {
		log.Errorln("Discord", err)
	}

	err = d.Open()
	if err != nil {
		log.Errorln("Discord", err)
	}

	d.AddHandler(commands.Entry)

	_, err = twitter.NewClient()
	fbOauth = facebook.NewClient()

	if err != nil {
		log.Error(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go webserver()
	go parallelism.Notify(rClient)
	go parallelism.GetProjects()
	go parallelism.GetProjectsCalendar()
	go parallelism.TickerHTTP(timer, commands.ProjectChan)
	go parallelism.Notifier(d, commands.ProjectChan)

	<-sc
	defer rClient.Close()
	defer d.Close()
}

func webserver() {
	var session *fb.Session
	var pageID, accessToken string

	db := redis.GetClient()
	authURL := fbOauth.AuthCodeURL(time.Now().String())

	log.Infof("Auth URL: %s", authURL)

	f := models.GetFacebook()
	token, err := db.Get("fbToken").Result()
	pageID, _ = db.Get("pageID").Result()

	if err == nil {
		session = facebook.GetClientByToken(token)
		session.SetAccessToken(token)
		f.PageID = pageID
		f.Session = session
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Receiving a request")

		t := time.Now()

		if err != nil {
			log.Warning("Token not found in Redis")
			code := r.FormValue("code")
			session, err = facebook.GetClient(fbOauth, code)

			if err != nil {
				log.WithFields(log.Fields{
					"time": t,
					"err":  err.Error(),
				}).Error("Facebook OAuth error")
			} else {
				if session != nil {

					r, err := session.Get("/me/accounts", fb.Params{
						"access_token": session.AccessToken(),
					})

					if err != nil {
						log.Error(err)
					} else {
						var items []fb.Result

						err := r.DecodeField("data", &items)

						if err != nil {
							log.Error(err)
						} else {

							pageID = items[0]["id"].(string)
							accessToken = items[0]["access_token"].(string)

							err = db.Set("fbToken", accessToken, 0).Err()
							if err != nil {
								log.Error(err)
							}

							err = db.Set("pageID", pageID, 0).Err()
							if err != nil {
								log.Error(err)
							}

							db.Save()
						}
					}
				}
			}
		}

		session.SetAccessToken(accessToken)

		f.PageID = pageID
		f.Session = session

		w.WriteHeader(http.StatusOK)
	})

	err = http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Error(err)
	}
}
