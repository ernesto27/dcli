package models

import "strings"

type Options struct {
	Cursor  int
	Choice  string
	Choices []string
}

func NewOptions() Options {
	return Options{
		Choices: []string{"Stop", "Start", "Remove"},
	}
}

func (o Options) View() string {
	s := strings.Builder{}
	s.WriteString("Options container [CONTAINER NAME]?\n\n")

	for i := 0; i < len(o.Choices); i++ {
		if o.Cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(o.Choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
