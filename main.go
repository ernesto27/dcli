package main

import (
	"context"
	"dockerniceui/docker"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

func createTable(title string, columns []string, rows [][]string) string {

	table := "# " + title + "\n\n"

	for _, column := range columns {
		table += "| " + column + " "
	}

	table += "|\n"

	table += "| ------------- | ------------- "

	if len(columns) > 2 {
		table += "| -------------"
	}

	table += "|\n"

	for _, row := range rows {
		for _, column := range row {
			table += "| " + column + " "
		}
		table += "|\n"
	}

	return table

}

func getContent(container docker.MyContainer) string {
	response := ""

	response += createTable("# Container status", []string{"Type", "Value"},
		[][]string{
			{"ID", container.ID},
			{"Name", container.Name},
			{"Image", container.Image},
			{"Status", container.State},
			{"Created", container.Status},
		})

	response += "\n\n---\n\n"
	rows := [][]string{}
	rows = append(rows, []string{"Command", container.Command})

	for _, env := range container.Env {
		rows = append(rows, []string{"ENV", env})
	}

	response += createTable("# Container detail", []string{"Type", "Value"}, rows)

	response += "\n\n---\n\n"
	response += createTable("# Networking", []string{"Network", "IP Address", "Gateway"}, [][]string{
		{container.Network.Name, container.Network.IPAddress, container.Network.Gateway},
	})

	return response

}

type myViewport struct {
	viewport viewport.Model
}

func newViewport(content string) (*myViewport, error) {
	const width = 120

	vp := viewport.New(width, 40)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &myViewport{
		viewport: vp,
	}, nil
}

func (vp myViewport) Init() tea.Cmd {
	return nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table      table.Model
	viewport   viewport.Model
	showDetail bool
}

func (m model) Init() tea.Cmd {
	// clean screen
	return tea.Batch(
		tea.ClearScreen,
		//tea.MoveTopLeft,
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
			// return m, tea.Batch(
			// 	tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			// )

			container, err := dockerClient.GetContainerByName(m.table.SelectedRow()[1])
			if err != nil {
				panic(err)
			}

			v, _ := newViewport(getContent(container))

			return model{
				table:      m.table,
				viewport:   v.viewport,
				showDetail: true,
			}, tea.ClearScreen

		}
	}

	m.table, _ = m.table.Update(msg)
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.showDetail {
		return m.viewport.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	}

	return baseStyle.Render(m.table.View()) + helpStyle("\n  ↑/↓: Navigate • Ctrl/C: Exit\n")
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

	m := model{
		table:    t,
		viewport: viewport.Model{},
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
