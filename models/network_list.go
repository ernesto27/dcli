package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
)

func NewNetworkList(rows []table.Row) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Name", Width: 40},
		{Title: "Driver", Width: 20},
		{Title: "IP Subnet", Width: 20},
		{Title: "IP Gateway", Width: 20},
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

func GetNetworkRows(networkList []types.NetworkResource, query string) []table.Row {
	var filtered []types.NetworkResource
	if query == "" {
		filtered = networkList
	} else {
		for _, network := range networkList {
			if strings.Contains(strings.ToLower(network.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(network.Driver), strings.ToLower(query)) {
				filtered = append(filtered, network)
			}
		}
	}

	var rows []table.Row
	for _, network := range filtered {
		subnet := ""
		gateway := ""
		if len(network.IPAM.Config) > 0 {
			subnet = network.IPAM.Config[0].Subnet
			gateway = network.IPAM.Config[0].Gateway
		}

		rows = append(rows, table.Row{
			network.ID,
			network.Name,
			network.Driver,
			subnet,
			gateway,
		})
	}

	return rows
}
