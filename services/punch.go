package services

import (
	"net/http"

	"golang.org/x/net/html"

	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/helpers"
	"github.com/wesleyholiveira/punchbot/models"
)

func GetProjects(endpoint string) []models.Project {

	projectMap := make(map[string]models.Project)
	projects := make([]models.Project, 0)

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode != 200 {
		log.Errorf("%s %d", resp.Status, resp.StatusCode)
	}

	sResponse := resp.Body
	defer resp.Body.Close()

	if err != nil {
		log.Error(err)
	}

	doc, err := html.Parse(sResponse)
	if err != nil {
		log.Error(err)
	}

	helpers.Transverse(doc, &projects, projectMap, "")

	/* re := regexp.MustCompile(`(\[.+\];)`)
	response := re.Find(sResponse)

	re = regexp.MustCompile(`\\\/`)
	response = re.ReplaceAll(response, []byte(`/`))
	response = response[:len(response)-1]

	err = json.Unmarshal(response, &projects)
	if err != nil {
		log.Error(err)
	} */

	return helpers.RemoveDuplicates(projects)
}
