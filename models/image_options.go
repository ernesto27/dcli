package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ImageOptions struct {
	Options
}

func NewImageOptions(image string) ImageOptions {
	choices := []string{Remove, ForceRemove}

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

			force := false
			if option == ForceRemove {
				force = true
			}

			err := m.dockerClient.ImageRemove(m.imageList.table.SelectedRow()[1], force)
			if err != nil {
				o.MessageError = err.Error()
				errAction = true
			}

			if !errAction {
				images, err := m.dockerClient.ImageList()
				if err != nil {
					o.MessageError = err.Error()
				}

				m.imageList = NewImageList(images, "")
				m.currentModel = MImageList
			}
		case "down":
			o.Cursor++
			if o.Cursor >= len(o.Choices) {
				o.Cursor = 0
			}

		case "up":
			o.Cursor--
			if o.Cursor < 0 {
				o.Cursor = len(o.Choices) - 1
			}

		}
	}

	return o, nil
}
