package parallelism

import (
	"math"
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
	projects := models.GetProjects()

	for t := range ticker.C {
		log.Infoln(t)

		current := services.GetProjects(endpoint, models.Home)
		if isNotEqualProjects(projects, &current) {
			*PrevProject = *projects
			*projects = current

			log.Info("Sending data to notifier")
			project <- projects
		}

	}
}

func isNotEqualProjects(prev *[]models.Project, current *[]models.Project) bool {
	currentVal := *current
	prevVal := *prev

	currentContent, _ := structhash.Hash(currentVal, 0)
	content, _ := structhash.Hash(prevVal, 0)

	log.Infof("prev: [%s](%d), current: [%s](%d)", content, len(prevVal), currentContent, len(currentVal))

	if currentContent != content && len(currentVal) > 0 {
		diff := int(math.Abs(float64(len(prevVal)) - float64(len(currentVal))))

		currentSlice := currentVal[0:diff]
		for _, c := range currentSlice {
			for _, p := range prevVal {
				if p.ID == c.ID {
					log.Warnf("Previosly project %s[%s] is equal to current project %s[%s]",
						p.Project, p.ID,
						c.Project, c.ID)
					log.Warn("IGNORED!!")
					return false
				}
			}
		}

		return true
	}

	return false
}
