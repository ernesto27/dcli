package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/volume"
)

type VolumeList struct {
	table table.Model
}

func NewVolumeList(volumeList []*volume.Volume, query string) VolumeList {
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Stack", Width: 20},
		{Title: "Driver", Width: 20},
		{Title: "Created", Width: 30},
	}

	rows := GetVolumeRows(volumeList, query)

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

	return VolumeList{table: t}
}

func (vl VolumeList) View(commands string, dockerVersion string) string {
	return baseStyle.Render(vl.table.View()) + helpStyle("\n DockerVersion: "+dockerVersion+" \n"+commands)
}

func (vl VolumeList) Update(msg tea.Msg, m *model) (table.Model, tea.Cmd) {
	if m.currentModel != MVolumeList {
		return vl.table, nil
	}

	vl.table, _ = vl.table.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
		}
	}

	return vl.table, nil
}

func GetVolumeRows(volumeList []*volume.Volume, query string) []table.Row {
	var filtered []*volume.Volume
	if query == "" {
		filtered = volumeList
	} else {
		for _, v := range volumeList {
			if strings.Contains(strings.ToLower(v.Name), strings.ToLower(query)) {
				filtered = append(filtered, v)
			}
		}
	}

	var rows []table.Row
	for _, v := range filtered {

		rows = append(rows, table.Row{
			v.Name,
			"",
			v.Driver,
			v.CreatedAt,
		})
	}

	return rows
}
