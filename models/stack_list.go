package models

import (
	"strings"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StackList struct {
	table table.Model
}

func NewStackList(stack []docker.MyStack, query string) StackList {
	columns := []table.Column{
		{Title: "Name", Width: 50},
		{Title: "Created", Width: 50},
	}

	rows := GetStackRows(stack, query)

	t := table.New(
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
	t.SetStyles(s)

	return StackList{table: t}
}

func (sl StackList) View(commands string, m *model) string {
	dockerStats := statsStyle.Render(m.getDockerStats())
	return baseStyle.Render(sl.table.View()) + helpStyle("\n"+dockerStats+"\n"+commands)
}

func (sl StackList) Update(msg tea.Msg, m *model) (table.Model, tea.Cmd) {
	if m.currentModel != MStackList {
		return sl.table, nil
	}

	sl.table, _ = sl.table.Update(msg)
	return sl.table, nil
}

func GetStackRows(stack []docker.MyStack, query string) []table.Row {
	var filtered []docker.MyStack
	if query == "" {
		filtered = stack
	} else {
		for _, s := range stack {
			if strings.Contains(strings.ToLower(s.Resource.Name), strings.ToLower(query)) {
				filtered = append(filtered, s)
			}
		}
	}

	var rows []table.Row
	for _, v := range filtered {
		rows = append(rows, table.Row{
			v.Resource.Name,
			v.Resource.Created.Format("2006-01-02 15:04:05"),
		})
	}

	return rows
}
