package main

import (
	"context"
	"dockerniceui/docker"
	"dockerniceui/models"
	"dockerniceui/utils"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type model struct {
	table      table.Model
	viewport   viewport.Model
	textinput  textinput.Model
	showDetail bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":

			return model{
				table:      m.table,
				viewport:   m.viewport,
				showDetail: false,
			}, tea.ClearScreen

		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			container, err := dockerClient.GetContainerByName(m.table.SelectedRow()[1])
			if err != nil {
				panic(err)
			}

			v, _ := models.NewViewport(utils.GetContent(container, utils.CreateTable))

			return model{
				table:      m.table,
				viewport:   v,
				showDetail: true,
			}, tea.ClearScreen

		}
	}

	m.table, _ = m.table.Update(msg)
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	if m.showDetail {
		return m.viewport.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	}

	return baseStyle.Render(m.table.View()) + helpStyle("\n  ↑/↓: Navigate • Ctrl/c: Exit • Ctrl/f: Search \n")
}

var dockerClient *docker.Docker

func main() {
	ctx := context.Background()
	var err error
	dockerClient, err = docker.New(ctx)
	if err != nil {
		panic(err)
	}

	containerList, err := dockerClient.ContainerList()
	if err != nil {
		panic(err)
	}

	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Container", Width: 40},
		{Title: "Image", Width: 40},
		{Title: "Port", Width: 40},
		{Title: "Status", Width: 30},
	}

	rowsItems := []table.Row{}

	for _, c := range containerList {
		port := ""
		if len(c.Ports) > 0 {
			port = fmt.Sprintf("http://%s:%d", "localhost", c.Ports[0].PublicPort)
		}

		up := "\u2191"
		greenUpArrow := "\033[32m" + up + "\033[0m"

		downArrow := "\u2193"
		redDownArrow := "\033[31m" + downArrow + "\033[0m"

		currState := redDownArrow + " " + c.State
		if c.State == "running" {
			currState = greenUpArrow + " " + c.State
		}

		item := []string{c.ID, c.Name, c.Image, port, currState}
		rowsItems = append(rowsItems, item)
	}

	rows := rowsItems
	t := models.NewTable(columns, rows)

	m := model{
		table: t,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
