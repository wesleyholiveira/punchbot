package parallelism

import (
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

func GetProjects() {
	projects := models.GetProjects()
	if len(*projects) == 0 {
		*projects = services.GetProjects(configs.Home, models.Home)
	}
}
