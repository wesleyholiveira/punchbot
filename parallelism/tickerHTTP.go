package parallelism

import (
	"time"

	"github.com/cnf/structhash"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

var PrevProject *[]models.Project

func init() {
	PrevProject = new([]models.Project)
}

func TickerHTTP(ticker *time.Ticker, project chan *[]models.Project) {
	var endpoint = configs.Home
	prjs := models.GetProjects()

	for t := range ticker.C {
		projects := services.GetProjects(endpoint, models.Home)
		currentContent, _ := structhash.Hash(projects, 0)
		content, _ := structhash.Hash(*prjs, 0)

		log.Infof("prev: [%s](%d), current: [%s](%d)", content, len(*prjs), currentContent, len(projects))

		if currentContent != content {
			*PrevProject = *prjs
			*prjs = projects
			project <- prjs
		}

		log.Info(t)
	}
}
