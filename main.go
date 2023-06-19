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
	pager      viewport.Model
	showDetail bool
	showSearch bool
	showLogs   bool
	ready      bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
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
					pager:      m.pager,
					viewport:   v,
					showDetail: true,
					showLogs:   false,
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

		case "ctrl+l":
			containerLogs, err := dockerClient.ContainerLogs("d2d611f2646f")
			if err != nil {
				panic(err)
			}
			for _, l := range containerLogs {
				logs += l + "\n"
			}

			headerHeight := lipgloss.Height(models.HeaderView(m.pager, ""))
			p := viewport.New(widthScreen, heigthScreen)

			p.YPosition = headerHeight + 1
			p.SetContent(logs)

			return model{
				table:      m.table,
				viewport:   m.viewport,
				textinput:  m.textinput,
				pager:      p,
				showDetail: false,
				showSearch: false,
				showLogs:   true,
			}, tea.ClearScreen

		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(models.HeaderView(m.pager, ""))
		footerHeight := lipgloss.Height(models.FooterView(m.pager))
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			widthScreen = msg.Width
			heigthScreen = msg.Height - verticalMarginHeight
			m.pager.YPosition = headerHeight
			m.ready = true
			m.pager.YPosition = headerHeight + 1
		} else {
			m.pager.Width = msg.Width
			m.pager.Height = msg.Height - verticalMarginHeight
		}

	}

	m.table, _ = m.table.Update(msg)
	m.viewport, _ = m.viewport.Update(msg)
	m.textinput, _ = m.textinput.Update(msg)
	m.pager, cmd = m.pager.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	if m.showSearch {
		return fmt.Sprintf(
			"Search container by name\n\n%s\n\n%s",
			m.textinput.View(),
			"(esc to back)",
		) + "\n"
	}

	if m.showDetail {
		return m.viewport.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	}

	if m.showLogs {
		return fmt.Sprintf("%s\n%s\n%s", models.HeaderView(m.pager, "container name - image name"), m.pager.View(), models.FooterView(m.pager))
	}

	return baseStyle.Render(m.table.View()) + helpStyle("\n  ↑/↓: Navigate • Ctrl/c: Exit • Ctrl/f: Search \n")
}

var dockerClient *docker.Docker
var containerList []docker.MyContainer
var logs string
var widthScreen int
var heigthScreen int

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

	if _, err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
