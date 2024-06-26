package models

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ContainerOptions struct {
	Options
}

func NewContainerOptions(container string, image string) ContainerOptions {
	choices := []string{Stop, Start, Remove, Restart, Pause, Unpause}

	return ContainerOptions{
		Options{
			Cursor:  0,
			Choice:  "",
			Choices: choices,
			Text1:   container,
			Text2:   image,
		},
	}
}

func (o ContainerOptions) View() string {
	title := fmt.Sprintf("Options container: %s - %s", o.Text1, o.Text2)
	return o.Options.View(title)
}

func (o ContainerOptions) Update(msg tea.Msg, m *model) (ContainerOptions, tea.Cmd) {
	if m.currentModel != MContainerOptions {
		return o, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.currentModel != MContainerOptions {
				return o, nil
			}

			errAction := false
			switch m.containerOptions.Choices[m.containerOptions.Cursor] {
			case Stop:
				err := m.dockerClient.ContainerStop(m.ContainerID)
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
				time.Sleep(1 * time.Second)
			case Start:
				err := m.dockerClient.ContainerStart(m.ContainerID)
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			case Remove:
				err := m.dockerClient.ContainerRemove(m.ContainerID)
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			case Restart:
				err := m.dockerClient.ContainerRestart(m.ContainerID)
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			case Pause:
				err := m.dockerClient.ContainerPause(m.ContainerID)
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			case Unpause:
				err := m.dockerClient.ContainerUnpause(m.ContainerID)
				if err != nil {
					fmt.Println(err)
					errAction = true
				}
			}

			if !errAction {
				m.setContainerList()
				m.currentModel = MContainerList
				return o, tea.ClearScreen
			}

		}
	}

	return o, nil
}
