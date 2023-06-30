package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ImageOptions struct {
	Options
}

func NewImageOptions(container string, image string) ImageOptions {
	choices := []string{Remove}

	return ImageOptions{
		Options{
			Cursor:    0,
			Choice:    "",
			Choices:   choices,
			Container: container,
			Image:     image,
		},
	}
}

func (o ImageOptions) View() string {
	s := strings.Builder{}

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#3259A8")).
		Padding(1).
		MarginTop(1).
		MarginBottom(1)

	s.WriteString("\n")

	for i := 0; i < len(o.Choices); i++ {
		if o.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(o.Choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press Esc to go back)\n")

	options := fmt.Sprintf("Options container: %s - %s", o.Container, o.Image)
	if o.Container == "" {
		options = fmt.Sprintf("Options image: %s", o.Image)
	}

	return style.Render(options) + s.String()
}

func (o ImageOptions) Update(msg tea.Msg, m *model) (ImageOptions, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.currentModel != MImageOptions {
				return o, nil
			}

			errAction := false
			option := m.imageOptions.Choices[m.imageOptions.Cursor]

			if option == Remove {
				err := m.dockerClient.ImageRemove(m.imageList.table.SelectedRow()[1])
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			}

			if !errAction {
				images, err := m.dockerClient.ImageList()
				if err != nil {
					fmt.Println(err)
				}

				m.imageList = NewImageList(images, "")
				m.currentModel = MImageList
			}

		}
	}

	return o, nil
}
