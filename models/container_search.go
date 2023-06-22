package models

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func NewTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ti
}
