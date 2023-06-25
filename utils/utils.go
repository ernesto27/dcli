package utils

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/image"
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

func ReverseLines(lines []string) {
	i := 0
	j := len(lines) - 1

	for i < j {
		lines[i], lines[j] = lines[j], lines[i]
		i++
		j--
	}
}

func ReverseSlice[T image.HistoryResponseItem](s []T) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
}

func TrimValue(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
