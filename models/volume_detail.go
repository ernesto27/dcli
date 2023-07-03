package models

import (
	"dockerniceui/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/volume"
)

func NewVolumeDetail(volume *volume.Volume, createTable utils.CreateTableFunc) (viewport.Model, error) {
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

func getContentVolume(volume *volume.Volume) string {
	// 	content := `
	// # Volume detail

	// | Type | Value |
	// | ---- | ----- |
	// | ID | 1234567890 |
	// | Created | my-network |
	// | Mount path | bridge |
	// | Driver | local |
	// | Labels | false

	// # Containers using this volume

	// | Name | Mounted at | Read -only
	// | -- | ---- | ------------ | ------------ |
	// | 1234567890 | my-container |
	// | 1234567890 | my-container |

	//`
	response := ""
	response += utils.CreateTable("# Volume detail", []string{"Type", "Value"},
		[][]string{
			{"ID", volume.Name},
			{"Created", volume.CreatedAt},
			{"Mount path", volume.Mountpoint},
			{"Driver", volume.Driver},
			{"Labels", fmt.Sprintf("%v", volume.Labels)},
		})

	return response
}
