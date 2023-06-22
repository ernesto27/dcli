package models

import (
	"dockerniceui/docker"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	running = "running"
	exited  = "exited"
)

func NewContainerList(rows []table.Row) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Container", Width: 30},
		{Title: "Image", Width: 30},
		{Title: "Port", Width: 20},
		{Title: "Size", Width: 20},
		{Title: "Status", Width: 20},
	}

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

	return t
}

func GetContainerRows(containerList []docker.MyContainer, query string) []table.Row {
	var filtered []docker.MyContainer
	if query == "" {
		filtered = containerList
	} else {
		for _, container := range containerList {
			if strings.Contains(strings.ToLower(container.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(container.Image), strings.ToLower(query)) {
				filtered = append(filtered, container)
			}
		}
	}

	rowsItems := []table.Row{}

	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].State > filtered[j].State
	})

	for _, c := range filtered {
		port := ""
		if len(c.Ports) > 0 {
			port = fmt.Sprintf("http://%s:%d", "localhost", c.Ports[0].PublicPort)
		}

		up := "\u2191"
		greenUpArrow := "\033[32m" + up + "\033[0m"

		downArrow := "\u2193"
		redDownArrow := "\033[31m" + downArrow + "\033[0m"

		currState := redDownArrow + " " + c.State
		if c.State == running {
			currState = greenUpArrow + " " + c.State
		}

		item := []string{c.ID, c.Name, c.Image, port, c.Size, currState}
		rowsItems = append(rowsItems, item)
	}

	return rowsItems
}
