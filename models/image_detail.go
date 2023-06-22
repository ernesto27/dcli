package models

import (
	"dockerniceui/docker"
	"dockerniceui/utils"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func NewImageDetail(image docker.MyImage, createTable utils.CreateTableFunc) (viewport.Model, error) {
	content := getContentDetailImage(image)
	const width = 120

	vp := viewport.New(width, 30)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return viewport.Model{}, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return viewport.Model{}, err
	}

	vp.SetContent(str)

	return vp, nil
}

func getContentDetailImage(image docker.MyImage) string {
	response := ""

	response += utils.CreateTable("# Image detail", []string{"Type", "Value"},
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

	response += utils.CreateTable("# Dockerfile details", []string{"Type", "Value"},
		[][]string{
			{"Author", image.Inspect.Author},
			{"CMD", cmd},
			{"Ports", fmt.Sprintf("%v", image.Inspect.Config.ExposedPorts)},
			{"Envs", fmt.Sprintf("%v", image.Inspect.Config.Env)},
		})

	response += "\n\n---\n\n"
	response += "# Docker hub image url \n " + utils.GetDockerHubURL(image.Summary.RepoTags[0])

	return response
}
