package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Options struct {
	Cursor    int
	Choice    string
	Choices   []string
	Container string
	Image     string
}

const (
	Stop    = "Stop"
	Start   = "Start"
	Remove  = "Remove"
	Restart = "Restart"
)

func NewContainerOptions(container string, image string) Options {
	var choices []string
	if container != "" {
		choices = []string{Stop, Start, Remove, Restart}
	} else {
		choices = []string{Remove}
	}

	return Options{
		Choices:   choices,
		Container: container,
		Image:     image,
	}
}

func (o Options) View() string {
	s := strings.Builder{}

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#3259A8")).
		Padding(1).
		MarginTop(1).
		MarginBottom(1)

	s.WriteString("\n")

	for i := 0; i < len(o.Choices); i++ {
		if o.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(o.Choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press Esc to go back)\n")

	options := fmt.Sprintf("Options container: %s - %s", o.Container, o.Image)
	if o.Container == "" {
		options = fmt.Sprintf("Options image: %s", o.Image)
	}

	return style.Render(options) + s.String()
}

func (o Options) Update(msg tea.Msg, m *model) (Options, tea.Cmd) {
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
