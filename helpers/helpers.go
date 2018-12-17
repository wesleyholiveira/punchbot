package helpers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/cnf/structhash"
	"github.com/wesleyholiveira/punchbot/models"
	"golang.org/x/net/html"
)

func RemoveDuplicates(elements []models.Project) []models.Project {
	encountered := map[string]bool{}
	var result []models.Project

	for _, v := range elements {
		if !encountered[v.Project] {
			encountered[v.Project] = true
			result = append(result, v)
		}
	}
	return result
}

func TransverseCalendar(n *html.Node, projects *[]models.Project, projectMap map[string]models.Project, key, day string) {

	project := new(models.Project)

	if n.Type == html.ElementNode {

		if n.Data == "li" {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					if attr.Val == "events-group" {
						day = n.FirstChild.FirstChild.LastChild.Data
					}
				}

				if attr.Key == "data-id" {
					key = attr.Val
					project.IDProject = attr.Val
					projectMap[key] = *project
				}
			}
		}

		if n.Data == "em" {
			prj := projectMap[key]
			prj.Project = n.FirstChild.Data
			prj.Day = day
			*projects = append(*projects, prj)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TransverseCalendar(c, projects, projectMap, key, day)
	}
}

func Transverse(n *html.Node, project *models.Project) {
	if n.Type == html.ElementNode {

		if n.Data == "p" {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					if strings.Contains(attr.Val, "projeto-box-sinopse") {
						if data := n.FirstChild; data != nil {
							project.Description = data.Data
							break
						}
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Transverse(c, project)
	}
}

func JsonUpdateToStruct(body io.Reader, projects *[]models.Project) error {
	r, _ := ioutil.ReadAll(body)
	re := regexp.MustCompile(`(\[.+\];)`)
	response := re.Find(r)

	re = regexp.MustCompile(`\\\/`)
	response = re.ReplaceAll(response, []byte(`/`))
	response = response[:len(response)-1]

	err := json.Unmarshal(response, projects)

	if err == nil {
		for i, project := range *projects {
			hash, _ := structhash.Hash(project.Project+project.Number, 0)
			(*projects)[i].HashID = hash
			(*projects)[i].AlreadyReleased = false
		}
		return nil
	}

	return err
}

func ParseChannels(channels string) map[string]string {
	reChannel := regexp.MustCompile(`\D+`)
	reTags := regexp.MustCompile(`\[(.*)\]`)
	mChannels := make(map[string]string)
	aChannels := strings.Split(channels, ",")
	for _, channel := range aChannels {
		channelID := reChannel.ReplaceAllString(channel, "")
		tags := reTags.FindAllStringSubmatch(channel, len(channel))
		tag := ""
		if len(tags) > 0 {
			tag = tags[0][1]
		}
		mChannels[channelID] = tag
	}
	return mChannels
}
