package parallelism

import (
	"time"

	"github.com/cnf/structhash"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

var firstTime bool

func init() {
	firstTime = true
}

func TickerHTTP(ticker *time.Ticker, project chan *[]models.Project) {
	prjs := models.GetProjects()
	var endpoint = configs.Home

	if firstTime {
		endpoint = configs.Calendar
		firstTime = false
	}

	for t := range ticker.C {
		projects := services.GetProjects(endpoint)
		currentContent, _ := structhash.Hash(projects, 0)
		content, _ := structhash.Hash(*prjs, 0)

		if currentContent != content {
			*prjs = projects
			project <- prjs
		}

		log.Info(t)
	}
}
