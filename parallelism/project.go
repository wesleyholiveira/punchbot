package parallelism

import (
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

func GetProjects() {
	projects := models.GetProjects()
	if len(*projects) == 0 {
		log.Info("Getting the latest releases")
		*projects = services.GetProjects(configs.Home, models.Home)
	}
}
