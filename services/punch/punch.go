package punch

import (
	"errors"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/helpers"
	"github.com/wesleyholiveira/punchbot/models"
	"golang.org/x/net/html"
)

var client *http.Client

func init() {
	url, _ := url.Parse(configs.ProxyURL)
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(url),
		},
	}
}

func Get(endpoint string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode > 399 {
		return nil, errors.New(resp.Status)
	}

	return resp, nil
}

func GetProjects(endpoint string, typ models.GetProjectsType) ([]models.Project, error) {
	projectMap := make(map[string]models.Project)
	projects := make([]models.Project, 0)

	sResponse, err := Get(endpoint)

	if err != nil {
		return nil, err
	} else {
		body := sResponse.Body
		defer body.Close()

		if typ == models.Calendar {
			doc, err := html.Parse(body)
			if err != nil {
				return nil, err
			}

			helpers.TransverseCalendar(doc, &projects, projectMap, "", "")
		} else {
			err = helpers.JsonUpdateToStruct(body, &projects)
			if err != nil {
				return nil, err
			}
		}
	}
	return helpers.RemoveDuplicates(projects), nil
}
