package parallelism

import (
	"time"

	"github.com/cnf/structhash"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services/punch"
)

var change bool
var backup *[]models.Project

func init() {
	change = false
	backup = new([]models.Project)
}

func TickerHTTP(ticker *time.Ticker, project chan *[]models.Project) {
	var endpoint = configs.Home

	log.Info("Starting the timer")

	for t := range ticker.C {
		log.Infoln(t)

		prev := backup

		if len(*prev) == 0 {
			prev = models.GetProjects()
		}

		current := punch.GetProjects(endpoint, models.Home)
		if change {
			changeAllAlreadyRelased(&current, prev)
		}

		current = GetExtraInfos(&current)

		for i, p := range (*prev)[:2] {
			c := current[i]
			if p.ID == c.ID {
				current[i].Description = p.Description
			}
		}

		if isNotEqualProjects(prev, &current) {
			log.Info("Sending data to notifier")
			project <- &current
			*backup = current
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
		log.Info("ExtraInfos Length: ", len(currentVal[0].ExtraInfos))
		fhd := false
		for _, extra := range currentVal[0].ExtraInfos {
			if extra.Format == "fullhd" {
				fhd = true
				change = true
				break
			}
		}

		if !fhd {
			if len(currentVal[0].ExtraInfos) > 2 {
				change = true
			}
		}

		return true
	}

	change = false
	return false
}

func changeAllAlreadyRelased(current *[]models.Project, prev *[]models.Project) {
	for i, c := range *current {
		for _, p := range *prev {
			if c.ID == p.ID {
				(*current)[i].AlreadyReleased = true
			}
		}
	}
}
