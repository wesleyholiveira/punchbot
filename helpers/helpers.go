package helpers

import (
	"fmt"
	"regexp"
	"strings"

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

func ParseChannels(channels string) map[string]string {
	reChannel := regexp.MustCompile(`(\D+)?(\[.*\])`)
	reTags := regexp.MustCompile(`\[(.*)\]`)
	mChannels := make(map[string]string)
	aChannels := strings.Split(channels, ",")
	for _, channel := range aChannels {
		channelID := reChannel.ReplaceAllString(channel, "")
		fmt.Println(channel, channelID)
		tags := reTags.FindAllStringSubmatch(channel, len(channel))[0][1]
		mChannels[channelID] = tags
	}
	return mChannels
}
