package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Search struct {
	textInput textinput.Model
}

func NewSearch() Search {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Search{
		textInput: ti,
	}
}

func (s Search) View() string {
	return fmt.Sprintf(
		"Search container by name\n\n%s\n\n%s",
		s.textInput.View(),
		"(esc to back)",
	) + "\n"
}

func (s Search) Update(msg tea.Msg, m *model) (textinput.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.currentModel != MContainerSearch {
				return s.textInput, nil
			}
			value := m.containerSearch.textInput.Value()
			t := NewContainerList(GetContainerRows(m.dockerClient.Containers, value))
			m.containerList = t
			m.currentModel = MContainerList
		}
	}

	s.textInput, _ = s.textInput.Update(msg)
	return s.textInput, nil
}
