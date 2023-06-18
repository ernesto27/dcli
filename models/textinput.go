package models

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type TextInput struct {
	textInput textinput.Model
	err       error
}

func NewTextInpu() TextInput {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return TextInput{
		textInput: ti,
		err:       nil,
	}
}
