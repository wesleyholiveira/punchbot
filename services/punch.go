package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/helpers"
	"github.com/wesleyholiveira/punchbot/models"
)

func GetProjects(endpoint string) []models.Project {
	var projects []models.Project

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode != 200 {
		log.Errorf("%s %d", resp.Status, resp.StatusCode)
	}

	sResponse, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Error(err)
	}

	re := regexp.MustCompile(`(\[.+\];)`)
	response := re.Find(sResponse)

	re = regexp.MustCompile(`\\\/`)
	response = re.ReplaceAll(response, []byte(`/`))
	response = response[:len(response)-1]

	err = json.Unmarshal(response, &projects)
	if err != nil {
		log.Error(err)
	}

	return helpers.RemoveDuplicates(projects)
}
