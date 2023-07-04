package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type VolumeOptions struct {
	Options
}

func NewVolumeOptions(name string) VolumeOptions {
	choices := []string{Remove}

	return VolumeOptions{
		Options{
			Cursor:  0,
			Choice:  "",
			Choices: choices,
			Text1:   name,
		},
	}
}

func (v VolumeOptions) View() string {
	title := fmt.Sprintf("Options volume: %s", v.Text1)
	return v.Options.View(title)
}

func (v VolumeOptions) Update(msg tea.Msg, m *model) (VolumeOptions, tea.Cmd) {
	if m.currentModel != MVolumeOptions {
		return v, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			errAction := false
			option := m.volumeOptions.Choices[m.volumeOptions.Cursor]

			if option == Remove {
				err := m.dockerClient.VolumeRemove(m.volumeList.table.SelectedRow()[0])
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			}

			if !errAction {
				v, err := m.dockerClient.VolumeList()
				if err != nil {
					fmt.Println(err)
				}

				m.volumeList = NewVolumeList(v, "")
				m.currentModel = MVolumeList
			}

		}
	}

	return v, nil
}
