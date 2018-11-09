package models

type Project struct {
	ID              string `json:"id"`
	IDProject       string `json:"id_projeto"`
	Project         string `json:"projeto"`
	Link            string `json:"link"`
	Number          string `json:"numero"`
	Screen          string `json:"screen"`
	Day             string
	HashID          string
	AlreadyReleased bool
}

var projects, calendarProjects *[]Project

func init() {
	projects = new([]Project)
	calendarProjects = new([]Project)
}

func GetProjects() *[]Project {
	return projects
}

func GetCalendarProjects() *[]Project {
	return calendarProjects
}
