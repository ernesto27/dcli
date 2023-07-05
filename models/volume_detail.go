package models

import (
	"dcli/docker"
	"dcli/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func NewVolumeDetail(volume docker.MyVolume, createTable utils.CreateTableFunc) (viewport.Model, error) {
	content := getContentVolume(volume)
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

func getContentVolume(v docker.MyVolume) string {
	response := ""
	response += utils.CreateTable("# Volume detail", []string{"Type", "Value"},
		[][]string{
			{"ID", v.Volume.Name},
			{"Created", v.Volume.CreatedAt},
			{"Mount path", v.Volume.Mountpoint},
			{"Driver", v.Volume.Driver},
			{"Labels", fmt.Sprintf("%v", v.Volume.Labels)},
		})

	response += "\n\n"

	if len(v.Containers) > 0 {
		rows := [][]string{}
		for _, c := range v.Containers {
			ro := "false"
			if c.ReadOnly {
				ro = "true"
			}
			rows = append(rows, []string{c.Name, c.MountedAt, ro})
		}

		response += utils.CreateTable("# Containers using this volume", []string{"Name", "Mounted at", "Read -only"}, rows)

	}

	return response
}
