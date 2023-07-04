package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ImageOptions struct {
	Options
}

func NewImageOptions(image string) ImageOptions {
	choices := []string{Remove}

	return ImageOptions{
		Options{
			Cursor:  0,
			Choice:  "",
			Choices: choices,
			Text1:   image,
		},
	}
}

func (o ImageOptions) View() string {
	title := fmt.Sprintf("Options image: %s", o.Text1)
	return o.Options.View(title)
}

func (o ImageOptions) Update(msg tea.Msg, m *model) (ImageOptions, tea.Cmd) {
	if m.currentModel != MImageOptions {
		return o, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
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
