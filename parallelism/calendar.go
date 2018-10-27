package parallelism

import (
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services"
)

func GetProjectsCalendar() {
	projectsCalendar := models.GetCalendarProjects()
	if len(*projectsCalendar) == 0 {
		*projectsCalendar = services.GetProjects(configs.Calendar, models.Calendar)
	}
}
