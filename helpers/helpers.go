package helpers

import (
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

func Transverse(n *html.Node, projects *[]models.Project, projectMap map[string]models.Project, key string) {

	project := new(models.Project)

	if n.Type == html.ElementNode {

		if n.Data == "li" {
			for _, attr := range n.Attr {
				if attr.Key == "data-id" {
					key = attr.Key
					project.IDProject = attr.Val
					projectMap[key] = *project
					break
				}
			}
		}

		if n.Data == "em" {
			prj := projectMap[key]
			prj.Project = n.FirstChild.Data
			*projects = append(*projects, prj)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Transverse(c, projects, projectMap, key)
	}
}
