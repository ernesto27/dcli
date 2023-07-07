package models

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ernesto27/dcli/utils"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	running = "running"
	exited  = "exited"
)

type ContainerList struct {
	table table.Model
}

func NewContainerList(rows []table.Row) ContainerList {
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

	return ContainerList{table: t}
}

func (cl ContainerList) View(commands string, m *model) string {
	dockerStats := statsStyle.Render(m.getDockerStats())
	return baseStyle.Render(cl.table.View()) + helpStyle("\n"+dockerStats+"\n"+commands)
}

func (cl ContainerList) Update(msg tea.Msg, m *model) (table.Model, tea.Cmd) {
	cl.table, _ = cl.table.Update(msg)
	if m.currentModel != MContainerList {
		return cl.table, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			container, err := m.dockerClient.GetContainerByName(m.containerList.table.SelectedRow()[1])
			if err != nil {
				fmt.Println(err)
			}
			vp, err := NewContainerDetail(container, utils.CreateTable)
			if err != nil {
				fmt.Println(err)
			}

			m.containerDetail = vp
			m.currentModel = MContainerDetail
			return cl.table, nil
		case "ctrl+f":
			m.containerSearch.textInput.SetValue("")
			m.currentModel = MContainerSearch
		case "ctrl+o":
			ov := NewContainerOptions(m.containerList.table.SelectedRow()[1], m.containerList.table.SelectedRow()[2])
			m.containerOptions = ov
			m.currentModel = MContainerOptions
			m.ContainerID = m.containerList.table.SelectedRow()[0]
		case "ctrl+l":
			containerLogs, err := m.dockerClient.ContainerLogs(m.containerList.table.SelectedRow()[0])
			if err != nil {
				panic(err)
			}

			headerHeight := lipgloss.Height(HeaderView(m.containerLogs.pager, m.containerList.table.SelectedRow()[1]))
			lv := NewContainerLogs(m.widthScreen, m.heightScreen, containerLogs, headerHeight)
			lv.container = m.containerList.table.SelectedRow()[1]
			lv.image = m.containerList.table.SelectedRow()[2]
			m.containerLogs = lv
			m.currentModel = MContainerLogs
		case "ctrl+b":
			images, err := m.dockerClient.ImageList()
			if err != nil {
				fmt.Println(err)
			}

			m.imageList = NewImageList(images, "")
			m.currentModel = MImageList
		case "ctrl+n":
			networks, err := m.dockerClient.NetworkList()
			if err != nil {
				fmt.Println(err)
			}

			m.networkList = NewNetworkList(networks, "")
			m.currentModel = MNetworkList
		case "ctrl+s":
			stats, err := m.dockerClient.ContainerStats(m.containerList.table.SelectedRow()[0])
			if err != nil {
				fmt.Println(err)
			}

			cs, err := NewContainerStats(stats, utils.CreateTable)
			if err != nil {
				fmt.Println(err)
			}
			m.containerStats = cs
			m.currentModel = MContainerStats
		}
	}

	return cl.table, nil
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
