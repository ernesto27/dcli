package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type VolumeSearch struct {
	Search
}

func NewVolumeSearch() VolumeSearch {
	return VolumeSearch{
		NewSearch(),
	}
}

func (vs VolumeSearch) View() string {
	return fmt.Sprintf(
		"Search volume by name\n\n%s\n\n%s",
		vs.textInput.View(),
		"(esc to back)",
	) + "\n"
}

func (vs VolumeSearch) Update(msg tea.Msg, m *model) (VolumeSearch, tea.Cmd) {
	if m.currentModel != MVolumeSearch {
		return vs, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			value := vs.textInput.Value()
			m.volumeList = NewVolumeList(m.dockerClient.Volumes, value)
			m.currentModel = MVolumeList
		}
	}

	vs.textInput, _ = vs.textInput.Update(msg)
	return vs, nil
}
