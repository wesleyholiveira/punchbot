package parallelism

import (
	"encoding/json"
	"io/ioutil"
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
		p, err := punch.GetProjects(configs.Home, models.Home)
		if err != nil {
			log.Error(err)
		} else {
			*projects = p
			aprojects := (*projects)[:2]
			*projects = GetExtraInfos(&aprojects)
			for _, el := range p {
				*projects = append(*projects, el)
			}
		}
	}
}

func GetExtraInfos(projects *[]models.Project) []models.Project {
	link := configs.PunchEndpoint + "/buscarVersoes"
	for i, p := range *projects {
		if !p.AlreadyReleased {
			log.Infof("Getting extra informations for %s project", p.Project)
			descEndpoint := configs.PunchEndpoint + p.Link
			values := url.Values{}

			req, err := punch.Get(descEndpoint)
			if err != nil {
				log.Errorf("Fail to get the page of project %s: ", p.Project)
				log.Error(err)
			} else {
				body := req.Body
				defer body.Close()

				doc, err := html.Parse(body)
				if err != nil {
					log.Error("Fail to parse HTML:", err)
				}

				helpers.Transverse(doc, &p)

				values.Set("id", p.ID)
				values.Set("projeto", p.IDProject)

				r, err := punch.HttpClient.PostForm(link, values)
				if err != nil {
					log.Error("Request error:", err)
				} else {
					body = r.Body
					defer body.Close()

					infos, _ := ioutil.ReadAll(body)
					extraInfos := &models.Infos{}
					err = json.Unmarshal(infos, extraInfos)
					if err != nil {
						log.Error("Parse error: ", err)
					} else {
						for _, itens := range extraInfos.Infos {
							p.ExtraInfos = append(p.ExtraInfos, itens)
						}

						(*projects)[i].Description = p.Description
						(*projects)[i].ExtraInfos = p.ExtraInfos
					}
				}
			}
		}
	}

	return *projects
}
