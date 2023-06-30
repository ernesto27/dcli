package models

import (
	"github.com/charmbracelet/bubbles/textinput"
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
