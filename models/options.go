package models

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Options struct {
	Cursor       int
	Choice       string
	Choices      []string
	Text1        string
	Text2        string
	MessageError string
}

const (
	Stop        = "Stop"
	Start       = "Start"
	Remove      = "Remove"
	ForceRemove = "Force Remove"
	Restart     = "Restart"
	Pause       = "Pause"
	Unpause     = "Unpause"
)

func (o Options) View(title string) string {
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
	return style.Render(title) + s.String() + "\n" + o.MessageError
}
