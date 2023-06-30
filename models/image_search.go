package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ImageSearch struct {
	Search
}

func NewImageSearch() ImageSearch {
	return ImageSearch{
		NewSearch(),
	}
}

func (is ImageSearch) View() string {
	return fmt.Sprintf(
		"Search image by name\n\n%s\n\n%s",
		is.textInput.View(),
		"(esc to back)",
	) + "\n"
}

func (is ImageSearch) Update(msg tea.Msg, m *model) (ImageSearch, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.currentModel != MImageSearch {
				return is, nil
			}
			value := m.imageSearch.textInput.Value()
			imgList := NewImageList(m.dockerClient.Images, value)
			m.imageList = imgList
			m.currentModel = MImageList
		}
	}

	is.textInput, _ = is.textInput.Update(msg)
	return is, nil
}
