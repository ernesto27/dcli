package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ContainerSearch struct {
	Search
}

func NewContainerSearch() ContainerSearch {
	return ContainerSearch{
		NewSearch(),
	}
}

func (cs ContainerSearch) View() string {
	return fmt.Sprintf(
		"Search container by name\n\n%s\n\n%s",
		cs.textInput.View(),
		"(esc to back)",
	) + "\n"
}

func (cs ContainerSearch) Update(msg tea.Msg, m *model) (ContainerSearch, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.currentModel != MContainerSearch {
				return cs, nil
			}
			value := cs.textInput.Value()
			t := NewContainerList(GetContainerRows(m.dockerClient.Containers, value))
			m.containerList = t
			m.currentModel = MContainerList
		}
	}

	cs.textInput, _ = cs.textInput.Update(msg)
	return cs, nil
}
