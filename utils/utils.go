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

func GetContentDetailImage(image docker.MyImage, createTable createTableFunc) string {
	response := ""

	response += createTable("# Image detail", []string{"Type", "Value"},
		[][]string{
			{"ID", image.Summary.ID},
			{"Name", image.Summary.RepoTags[0]},
			{"Size", image.GetFormatSize()},
			{"Created", image.GetFormatTimestamp()},
			{"Build", image.Inspect.Os + " - " + image.Inspect.Architecture + " - Docker version " + image.Inspect.DockerVersion},
		})

	cmd := ""
	if len(image.Inspect.Config.Cmd) > 0 {
		cmd = strings.Join(image.Inspect.Config.Cmd, " ")
	}

	response += createTable("# Dockerfile details", []string{"Type", "Value"},
		[][]string{
			{"Author", image.Inspect.Author},
			{"CMD", cmd},
			{"Ports", fmt.Sprintf("%v", image.Inspect.Config.ExposedPorts)},
			{"Envs", fmt.Sprintf("%v", image.Inspect.Config.Env)},
		})

	response += "\n\n---\n\n"
	response += "# Docker hub image url \n " + GetDockerHubURL(image.Summary.RepoTags[0])

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
