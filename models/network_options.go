package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type NetworkOptions struct {
	Options
}

func NewNetworkOptions(name string) NetworkOptions {
	choices := []string{Remove}

	return NetworkOptions{
		Options{
			Cursor:  0,
			Choice:  "",
			Choices: choices,
			Text1:   name,
		},
	}
}

func (n NetworkOptions) View() string {
	title := fmt.Sprintf("Options network: %s", n.Text1)
	return n.Options.View(title)
}

func (n NetworkOptions) Update(msg tea.Msg, m *model) (NetworkOptions, tea.Cmd) {
	if m.currentModel != MNetworkOptions {
		return n, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			errAction := false
			option := m.networkOptions.Choices[m.imageOptions.Cursor]

			if option == Remove {
				err := m.dockerClient.NetworkRemove(m.networkList.table.SelectedRow()[0])
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			}

			if !errAction {
				networks, err := m.dockerClient.NetworkList()
				if err != nil {
					fmt.Println(err)
				}

				m.networkList = NewNetworkList(networks, "")
				m.currentModel = MNetworkList
			}

		}
	}

	return n, nil
}
