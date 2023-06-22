package models

import (
	"fmt"
	"strings"
)

type Options struct {
	Cursor    int
	Choice    string
	Choices   []string
	Container string
	Image     string
}

const (
	Stop   = "Stop"
	Start  = "Start"
	Remove = "Remove"
)

func NewContainerOptions(container string, image string) Options {
	return Options{
		Choices:   []string{Stop, Start, Remove},
		Container: container,
		Image:     image,
	}
}

func (o Options) View() string {
	s := strings.Builder{}
	options := fmt.Sprintf("Options container: %s - %s \n\n", o.Container, o.Image)
	s.WriteString(options)

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

	return s.String()
}
