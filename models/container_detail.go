package models

import (
	"dockerniceui/docker"
	"dockerniceui/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ContianerDetail struct {
	viewport viewport.Model
}

func NewContainerDetail(container docker.MyContainer, createTable utils.CreateTableFunc) (ContianerDetail, error) {
	content := getContent(container)

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
		return ContianerDetail{}, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return ContianerDetail{}, err
	}

	vp.SetContent(str)

	return ContianerDetail{viewport: vp}, nil
}

func (cd ContianerDetail) View() string {
	return cd.viewport.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
}

func (cd ContianerDetail) Update(msg tea.Msg, m *model) (viewport.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.setContainerList()
		}
	}

	cd.viewport, _ = cd.viewport.Update(msg)
	return cd.viewport, nil
}

func getContent(container docker.MyContainer) string {
	response := ""

	response += utils.CreateTable("# Container status", []string{"Type", "Value"},
		[][]string{
			{"ID", container.ID},
			{"Name", container.Name},
			{"Image", container.Image},
			{"Status", container.State},
			{"Created", container.Status},
		})

	response += "\n\n---\n\n"
	rows := [][]string{}

	ports := ""
	for _, c := range container.Ports {
		ports += fmt.Sprintf("%s:%d->%d/%s ", c.IP, c.PublicPort, c.PrivatePort, c.Type)
	}
	rows = append(rows, []string{"Ports", ports})
	rows = append(rows, []string{"Command", container.Command})

	for _, env := range container.Env {
		rows = append(rows, []string{"ENV", env})
	}

	response += utils.CreateTable("# Container detail", []string{"Type", "Value"}, rows)

	response += "\n\n---\n\n"
	response += utils.CreateTable("# Networking", []string{"Network", "IP Address", "Gateway"}, [][]string{
		{container.Network.Name, container.Network.IPAddress, container.Network.Gateway},
	})

	response += "\n\n---\n\n"

	response += "# Docker hub image url \n " + utils.GetDockerHubURL(container.Image)

	return response
}
