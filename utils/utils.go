package utils

import (
	"dockerniceui/docker"
	"fmt"
	"strings"
)

type createTableFunc func(string, []string, [][]string) string

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

func GetContent(container docker.MyContainer, createTable createTableFunc) string {
	response := ""

	response += createTable("# Container status", []string{"Type", "Value"},
		[][]string{
			{"ID", container.ID},
			{"Name", container.Name},
			{"Image", container.Image},
			{"Status", container.State},
			{"Created", container.Status},
		})

	response += "\n\n---\n\n"
	rows := [][]string{}
	rows = append(rows, []string{"Command", container.Command})

	for _, env := range container.Env {
		rows = append(rows, []string{"ENV", env})
	}

	response += createTable("# Container detail", []string{"Type", "Value"}, rows)

	response += "\n\n---\n\n"
	response += createTable("# Networking", []string{"Network", "IP Address", "Gateway"}, [][]string{
		{container.Network.Name, container.Network.IPAddress, container.Network.Gateway},
	})

	response += "\n\n---\n\n"

	response += "# Docker hub image url \n " + GetDockerHubURL(container.Image)

	return response
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
