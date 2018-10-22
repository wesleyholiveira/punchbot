package parallelism

import (
	"fmt"
	"time"

	"github.com/cnf/structhash"
	config "github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

func TickerHTTP(ticker *time.Ticker, notify chan bool) {
	prjs := models.GetProjects()
	for t := range ticker.C {
		projects := services.GetProjects(config.PunchEndpoint)
		currentContent, _ := structhash.Hash(projects, 0)
		content, _ := structhash.Hash(*prjs, 0)

		if currentContent != content || len(*prjs) == 0 {
			*prjs = projects
			notify <- true
		} else {
			notify <- false
		}

		fmt.Println(t)
	}
}
