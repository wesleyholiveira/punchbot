package models

type Project struct {
	ID        string `json:"id"`
	IDProject string `json:"id_projeto"`
	Link      string `json:"link"`
	Project   string `json:"projeto"`
	Number    string `json:"numero"`
	Screen    string `json:"screen"`
}

var projects *[]Project

func init() {
	projects = new([]Project)
}

func GetProjects() *[]Project {
	return projects
}
