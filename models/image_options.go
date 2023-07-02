package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
