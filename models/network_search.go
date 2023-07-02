package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type NetworkSearch struct {
	Search
}

func NewNetworkSearch() NetworkSearch {
	return NetworkSearch{
		NewSearch(),
	}
}

func (ns NetworkSearch) View() string {
	return fmt.Sprintf(
		"Search network by name\n\n%s\n\n%s",
		ns.textInput.View(),
		"(esc to back)",
	) + "\n"
}

func (ns NetworkSearch) Update(msg tea.Msg, m *model) (NetworkSearch, tea.Cmd) {
	if m.currentModel != MNetworkSearch {
		return ns, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			value := ns.textInput.Value()
			m.networkList = NewNetworkList(m.dockerClient.Networks, value)
			m.currentModel = MNetworkList
		}
	}

	ns.textInput, _ = ns.textInput.Update(msg)
	return ns, nil
}
