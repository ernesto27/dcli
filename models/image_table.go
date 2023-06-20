package models

import (
	"dockerniceui/utils"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
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

func GetImageRows(images []types.ImageSummary, query string) []table.Row {
	rowsItems := []table.Row{}
	for _, i := range images {
		fmt.Println(i.Size)
		item := []string{i.ID, i.RepoTags[0], utils.FormatSize(i.Size), FormatTimestamp(i.Created)}
		rowsItems = append(rowsItems, item)
	}

	return rowsItems
}

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	duration := time.Since(t)

	months := int(duration.Hours() / (24 * 30))
	weeks := int(duration.Hours() / (24 * 7))
	days := int(duration.Hours() / 24)

	plural := "s"
	if months > 0 {
		if months == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d month%s ago", months, plural)
	} else if weeks > 0 {
		if weeks == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d week%s ago", weeks, plural)
	} else if days > 0 {
		if days == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d day%s ago", days, plural)
	} else {
		return "today"
	}
}
