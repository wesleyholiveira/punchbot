package services

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/helpers"
	"github.com/wesleyholiveira/punchbot/models"
	"golang.org/x/net/html"
)

func Get(endpoint string) (*http.Response, error) {
	resp, err := http.Get(endpoint)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode > 399 {
		return nil, errors.New(resp.Status)
	}

	return resp, nil
}

func GetProjects(endpoint string, typ models.GetProjectsType) []models.Project {

	projectMap := make(map[string]models.Project)
	projects := make([]models.Project, 0)

	sResponse, err := Get(endpoint)

	if err != nil {
		log.Error(err)
	} else {
		body := sResponse.Body
		defer body.Close()

		if typ == models.Calendar {
			doc, err := html.Parse(body)
			if err != nil {
				log.Error(err)
			}

			helpers.Transverse(doc, &projects, projectMap, "", "")
		} else {
			helpers.JsonUpdateToStruct(body, &projects)
		}
	}
	return helpers.RemoveDuplicates(projects)
}
