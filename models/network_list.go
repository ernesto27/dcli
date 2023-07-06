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

type NetworkList struct {
	table table.Model
}

func NewNetworkList(networkList []docker.MyNetwork, query string) NetworkList {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Name", Width: 40},
		{Title: "Driver", Width: 20},
		{Title: "IP Subnet", Width: 20},
		{Title: "IP Gateway", Width: 20},
	}

	rows := GetNetworkRows(networkList, query)

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

	return NetworkList{table: t}
}

func (nl NetworkList) View(commands string, dockerVersion string) string {
	return baseStyle.Render(nl.table.View()) + helpStyle("\n DockerVersion: "+dockerVersion+" \n"+commands)
}

func (cl NetworkList) Update(msg tea.Msg, m *model) (table.Model, tea.Cmd) {
	cl.table, _ = m.networkList.table.Update(msg)
	if m.currentModel != MNetworkList {
		return cl.table, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			network, err := m.dockerClient.GetNetworkByName(m.networkList.table.SelectedRow()[1])
			if err != nil {
				fmt.Println(err)
			}

			nd, err := NewNetworkDetail(network, utils.CreateTable)
			if err != nil {
				fmt.Println(err)
			}
			m.networkDetail = nd
			m.currentModel = MNetworkDetail
		case "ctrl+f":
			m.networkSearch.textInput.SetValue("")
			m.currentModel = MNetworkSearch
		case "ctrl+o":
			m.networkOptions = NewNetworkOptions(m.networkList.table.SelectedRow()[1])
			m.currentModel = MNetworkOptions
		}
	}

	return cl.table, nil
}

func GetNetworkRows(networkList []docker.MyNetwork, query string) []table.Row {
	var filtered []docker.MyNetwork
	if query == "" {
		filtered = networkList
	} else {
		for _, network := range networkList {
			if strings.Contains(strings.ToLower(network.Resource.Name), strings.ToLower(query)) {
				filtered = append(filtered, network)
			}
		}
	}

	var rows []table.Row
	for _, network := range filtered {

		rows = append(rows, table.Row{
			network.Resource.ID,
			network.Resource.Name,
			network.Resource.Driver,
			network.Subnet,
			network.Gateway,
		})
	}

	return rows
}
