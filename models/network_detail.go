package models

import (
	"github.com/ernesto27/dcli/utils"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func NewNetworkDetail(network docker.MyNetwork, containers []docker.MyContainer, createTable utils.CreateTableFunc) (viewport.Model, error) {
	content := getContentNetwork(network, containers)
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

func getContentNetwork(network docker.MyNetwork, containers []docker.MyContainer) string {
	response := ""

	attachable := "false"
	if network.Resource.Attachable {
		attachable = "true"
	}

	response += utils.CreateTable("# Network status", []string{"Type", "Value"},
		[][]string{
			{"ID", network.Resource.ID},
			{"Name", network.Resource.Name},
			{"Driver", network.Resource.Driver},
			{"Attachable", attachable},
			{"Subnet", network.Subnet},
			{"Gateway", network.Gateway},
		})

	if len(containers) > 0 {
		columns := []string{"Name", "IPv4 Address"}
		rows := [][]string{}

		for _, c := range containers {
			rows = append(rows, []string{c.Name, c.Network.IPAddress})
		}

		response += "\n\n"
		response += utils.CreateTable("# Containers", columns, rows)
	}

	return response
}
