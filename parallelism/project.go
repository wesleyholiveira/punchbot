package parallelism

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wesleyholiveira/punchbot/helpers"
	"golang.org/x/net/html"

	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services/punch"
)

func GetProjects() {
	projects := models.GetProjects()
	if len(*projects) == 0 {
		log.Info("Getting the latest releases")
		*projects = punch.GetProjects(configs.Home, models.Home)
		*projects = GetExtraInfos(projects)
	}
}

func GetExtraInfos(projects *[]models.Project) []models.Project {
	link := configs.PunchEndpoint + "/buscarVersoes"
	for i, p := range (*projects)[0:2] {
		log.Infof("Getting extra informations for %s project", p.Project)
		if p.ExtraInfos == nil {
			descEndpoint := configs.PunchEndpoint + p.Link
			values := url.Values{}

			req, err := punch.Get(descEndpoint)
			if err != nil {
				log.Error("Fail to get the page of project: ", err)
			}

			body := req.Body
			defer body.Close()

			doc, err := html.Parse(body)
			if err != nil {
				log.Error("Fail to parse HTML:", err)
			}

			helpers.Transverse(doc, &p)

			values.Set("id", p.ID)
			values.Set("projeto", p.IDProject)

			r, err := http.PostForm(link, values)
			if err != nil {
				log.Error("Request error:", err)
			}

			body = r.Body
			defer body.Close()

			if body != nil {
				infos, _ := ioutil.ReadAll(body)
				extraInfos := &models.Infos{}
				err := json.Unmarshal(infos, extraInfos)
				if err != nil {
					log.Error("Parse error: ", err)
				}

				for _, itens := range extraInfos.Infos {
					p.ExtraInfos = append(p.ExtraInfos, itens)
				}

				(*projects)[i].Description = p.Description
				(*projects)[i].ExtraInfos = p.ExtraInfos
			}
		}
	}

	return *projects
}
