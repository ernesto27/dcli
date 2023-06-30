package models

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
