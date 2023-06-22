package utils

import (
	"fmt"
	"strings"
)

type CreateTableFunc func(string, []string, [][]string) string

func CreateTable(title string, columns []string, rows [][]string) string {
	table := "# " + title + "\n\n"

	for _, column := range columns {
		table += "| " + column + " "
	}

	table += "|\n"

	table += "| ------------- | ------------- "

	if len(columns) > 2 {
		table += "| -------------"
	}

	table += "|\n"

	for _, row := range rows {
		for _, column := range row {
			table += "| " + column + " "
		}
		table += "|\n"
	}

	return table

}

func GetDockerHubURL(image string) string {
	imageParts := strings.Split(image, ":")
	image = ""
	pathDefault := "_"
	if len(imageParts) > 0 {
		image = imageParts[0]
	}

	if strings.Contains(image, "/") {
		pathDefault = "r"
	}

	dockerHubLink := fmt.Sprintf("https://hub.docker.com/%s/%s", pathDefault, image)

	return dockerHubLink
}
