package parallelism

import (
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services/punch"
)

func GetProjects() {
	projects := models.GetProjects()
	if len(*projects) == 0 {
		log.Info("Getting the latest releases")
		*projects = punch.GetProjects(configs.Home, models.Home)
	}
}
