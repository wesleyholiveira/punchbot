package parallelism

import (
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

func GetProjects() {
	log.Info("Getting the latest releases")
	projects := models.GetProjects()
	if len(*projects) == 0 {
		*projects = services.GetProjects(configs.Home, models.Home)
	}
}
