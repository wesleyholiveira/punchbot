package models

type Project struct {
	IDProject string `json:"id_projeto"`
	Project   string `json:"projeto"`
}

var projects *[]Project

func init() {
	projects = new([]Project)
}

func GetProjects() *[]Project {
	return projects
}
