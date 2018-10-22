package helpers

import "github.com/wesleyholiveira/punchbot/models"

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
