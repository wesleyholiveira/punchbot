package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/helpers"
	"github.com/wesleyholiveira/punchbot/models"
	"golang.org/x/net/html"
)

func Get(endpoint string) *http.Response {
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode != 200 {
		log.Errorf("%s %d", resp.Status, resp.StatusCode)
	}

	if err != nil {
		log.Error(err)
	}

	return resp
}

func GetProjects(endpoint string, typ models.GetProjectsType) []models.Project {

	projectMap := make(map[string]models.Project)
	projects := make([]models.Project, 0)

	sResponse := Get(endpoint).Body
	defer sResponse.Close()

	if typ == models.Calendar {
		doc, err := html.Parse(sResponse)
		if err != nil {
			log.Error(err)
		}

		helpers.Transverse(doc, &projects, projectMap, "")
	} else {
		r, _ := ioutil.ReadAll(sResponse)
		re := regexp.MustCompile(`(\[.+\];)`)
		response := re.Find(r)

		re = regexp.MustCompile(`\\\/`)
		response = re.ReplaceAll(response, []byte(`/`))
		response = response[:len(response)-1]

		err := json.Unmarshal(response, &projects)
		if err != nil {
			log.Error(err)
		}
	}

	return helpers.RemoveDuplicates(projects)
}
