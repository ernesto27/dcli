package models

import (
	"dockerniceui/docker"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func NewImageTable(rows []table.Row) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Image", Width: 40},
		{Title: "Size", Width: 20},
		{Title: "Created", Width: 20},
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

func GetImageRows(images []docker.MyImage, query string) []table.Row {
	rowsItems := []table.Row{}
	for _, i := range images {
		fmt.Println(i.Summary.Size)
		item := []string{i.Summary.ID, i.Summary.RepoTags[0], i.GetFormatSize(), i.GetFormatTimestamp()}
		rowsItems = append(rowsItems, item)
	}

	return rowsItems
}
