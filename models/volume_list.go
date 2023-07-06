package models

import (
	"fmt"
	"strings"

	"github.com/ernesto27/dcli/utils"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VolumeList struct {
	table table.Model
}

func NewVolumeList(volumeList []docker.MyVolume, query string) VolumeList {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Stack", Width: 10},
		{Title: "Driver", Width: 10},
		{Title: "Mount point", Width: 40},
		{Title: "Created", Width: 25},
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
			v, err := m.dockerClient.GetVolumeByName(vl.table.SelectedRow()[0])
			if err != nil {
				fmt.Println(err)
			}

			vd, err := NewVolumeDetail(v, utils.CreateTable)
			if err != nil {
				fmt.Println(err)
			}
			m.volumeDetail = vd
			m.currentModel = MVolumeDetail
		case "ctrl+f":
			m.volumeSearch.textInput.SetValue("")
			m.currentModel = MVolumeSearch
		case "ctrl+o":
			m.volumeOptions = NewVolumeOptions(vl.table.SelectedRow()[0])
			m.currentModel = MVolumeOptions
		}
	}

	return vl.table, nil
}

func GetVolumeRows(volumeList []docker.MyVolume, query string) []table.Row {
	var filtered []docker.MyVolume
	if query == "" {
		filtered = volumeList
	} else {
		for _, v := range volumeList {
			if strings.Contains(strings.ToLower(v.Volume.Name), strings.ToLower(query)) {
				filtered = append(filtered, v)
			}
		}
	}

	var rows []table.Row
	for _, v := range filtered {

		rows = append(rows, table.Row{
			v.Volume.Name,
			"",
			v.Volume.Driver,
			v.Volume.Mountpoint,
			v.Volume.CreatedAt,
		})
	}

	return rows
}
