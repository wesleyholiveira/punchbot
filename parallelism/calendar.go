package parallelism

import (
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services/punch"
)

func GetProjectsCalendar() {
	log.Info("Getting the projects from the calendars' page")
	projectsCalendar := models.GetCalendarProjects()
	if len(*projectsCalendar) == 0 {
		p, err := punch.GetProjects(configs.Calendar, models.Calendar)
		if err != nil {
			log.Error(err)
		}
		*projectsCalendar = p
	}
}
