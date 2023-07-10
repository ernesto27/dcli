package models

import (
	"github.com/ernesto27/dcli/utils"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func NewStackDetail(stack docker.MyStack, createTable utils.CreateTableFunc) (viewport.Model, error) {
	content := getContentStack(stack)
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

func getContentStack(stack docker.MyStack) string {
	response := ""

	response += "\n # Stack detail " + stack.Resource.Name + "\n\n"

	if len(stack.Containers) > 0 {
		columns := []string{"Name", "Status", "Image", "IP"}
		data := [][]string{}
		for _, container := range stack.Containers {
			data = append(data, []string{
				container.Name,
				container.Status,
				container.Image,
				container.Network.IPAddress,
			})
		}
		response += utils.CreateTable("# Containers", columns, data)
	}

	return response
}
