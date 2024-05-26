package models

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
)

type ContainerTop struct {
	Table table.Model
	title string
}

func NewContainerTop(data container.ContainerTopOKBody, title string) ContainerTop {
	columns := []table.Column{
		{Title: "UID", Width: 10},
		{Title: "PID", Width: 10},
		{Title: "PPID", Width: 10},
		{Title: "C", Width: 6},
		{Title: "STIME", Width: 10},
		{Title: "TTY", Width: 10},
		{Title: "TIME", Width: 10},
		{Title: "CMD", Width: 40},
	}

	rows := []table.Row{}
	for _, process := range data.Processes {
		rows = append(rows, table.Row(process))
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(180),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tableModel.SetStyles(s)

	return ContainerTop{
		Table: tableModel,
		title: title,
	}
}

func (containerTop ContainerTop) View() string {
	return titleTableStyle(containerTop.title) + "\n\n" + containerTop.Table.View()
}

func (containerTop ContainerTop) Update(msg tea.Msg, m *model) (ContainerTop, tea.Cmd) {
	containerTop.Table, _ = containerTop.Table.Update(msg)
	return containerTop, nil
}
