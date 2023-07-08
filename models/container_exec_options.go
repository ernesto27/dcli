package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ContainerExecOptions struct {
	Options
}

const (
	Bash = "bin/bash"
	Sh   = "bin/sh"
	Ash  = "bin/ash"
)

func NewContainerExecOptions(container string) ContainerExecOptions {
	choices := []string{Bash, Sh, Ash}

	return ContainerExecOptions{
		Options{
			Cursor:  0,
			Choice:  "",
			Choices: choices,
			Text1:   container,
		},
	}
}

func (o ContainerExecOptions) View() string {
	title := fmt.Sprintf("Exec command container: %s", o.Text1)
	return o.Options.View(title)
}

func (o ContainerExecOptions) Update(msg tea.Msg, m *model) (ContainerExecOptions, tea.Cmd) {
	if m.currentModel != MContainerExecOptions {
		return o, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
