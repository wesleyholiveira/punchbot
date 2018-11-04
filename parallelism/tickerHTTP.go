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
	prev := models.GetProjects()

	log.Info("Starting the timer")

	for t := range ticker.C {
		log.Infoln(t)

		current := services.GetProjects(endpoint, models.Home)
		if isNotEqualProjects(prev, &current) {
			*PrevProject = *prev
			*prev = current

			log.Info("Sending data to notifier")
			project <- &current
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

		if diff == 0 {
			diff = 1
		}

		currentSlice := currentVal[0:diff]
		prevSlice := prevVal[0:diff]

		for _, c := range currentSlice {
			for _, p := range prevSlice {
				if p.HashID == c.HashID {
					log.Warnf("Previosly project %s[%s] is equal to current project %s[%s]",
						p.Project, p.HashID,
						c.Project, c.HashID)
					log.Warn("IGNORED!!")
					return false
				}
			}
		}

		return true
	}

	return false
}
