package models

type Project struct {
	IDProject string `json:"id_projeto"`
	Project   string `json:"projeto"`
	Link      string `json:"link"`
	Numero    string `json:"numero"`
	Screen    string `json:"screen"`
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
