package models

import (
	"fmt"
	"strings"

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
