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
	showSearch bool
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
			t := models.NewTable(models.GetContainerColumns(), models.GetContainerRows(containerList, ""))
			return model{
				table:      t,
				viewport:   m.viewport,
				textinput:  m.textinput,
				showDetail: false,
			}, tea.ClearScreen

		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.showSearch {
				value := m.textinput.Value()
				t := models.NewTable(models.GetContainerColumns(), models.GetContainerRows(containerList, value))

				return model{
					table:      t,
					textinput:  m.textinput,
					showDetail: false,
					showSearch: false,
				}, tea.ClearScreen

			} else {
				container, err := dockerClient.GetContainerByName(m.table.SelectedRow()[1])
				if err != nil {
					panic(err)
				}

				v, _ := models.NewViewport(utils.GetContent(container, utils.CreateTable))

				return model{
					table:      m.table,
					textinput:  m.textinput,
					viewport:   v,
					showDetail: true,
					showSearch: false,
				}, tea.ClearScreen

			}

		case "ctrl+f":
			m.textinput.SetValue("")
			return model{
				table:      m.table,
				viewport:   m.viewport,
				textinput:  m.textinput,
				showDetail: false,
				showSearch: true,
			}, tea.ClearScreen
		}
	}

	m.table, _ = m.table.Update(msg)
	m.viewport, _ = m.viewport.Update(msg)
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	if m.showSearch {
		return fmt.Sprintf(
			"Search container by name\n\n%s\n\n%s",
			m.textinput.View(),
			"(esc to quit)",
		) + "\n"
	}

	if m.showDetail {
		return m.viewport.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	}

	return baseStyle.Render(m.table.View()) + helpStyle("\n  ↑/↓: Navigate • Ctrl/c: Exit • Ctrl/f: Search \n")
}

var dockerClient *docker.Docker
var containerList []docker.MyContainer

func main() {
	ctx := context.Background()
	var err error
	dockerClient, err = docker.New(ctx)
	if err != nil {
		panic(err)
	}

	containerList, err = dockerClient.ContainerList()
	if err != nil {
		panic(err)
	}

	t := models.NewTable(models.GetContainerColumns(), models.GetContainerRows(containerList, ""))

	m := model{
		table:     t,
		textinput: models.NewTextInput(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
