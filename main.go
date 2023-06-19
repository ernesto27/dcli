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

type currentView int

const (
	ListContainer currentView = iota
	DetailContainer
	SearchContainer
	LogsContainer
	OptionsContainer
)

type LogsView struct {
	pager     viewport.Model
	container string
	image     string
}

type model struct {
	table       table.Model
	viewport    viewport.Model
	textinput   textinput.Model
	logsView    LogsView
	optionsView models.Options
	ready       bool
	currentView currentView
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
				table:       t,
				viewport:    m.viewport,
				textinput:   m.textinput,
				currentView: ListContainer,
			}, tea.ClearScreen

		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.currentView == SearchContainer {
				value := m.textinput.Value()
				t := models.NewTable(models.GetContainerColumns(), models.GetContainerRows(containerList, value))

				return model{
					table:     t,
					textinput: m.textinput,
				}, tea.ClearScreen

			} else {
				container, err := dockerClient.GetContainerByName(m.table.SelectedRow()[1])
				if err != nil {
					panic(err)
				}

				v, _ := models.NewViewport(utils.GetContent(container, utils.CreateTable))

				return model{
					table:       m.table,
					textinput:   m.textinput,
					logsView:    m.logsView,
					viewport:    v,
					currentView: DetailContainer,
				}, tea.ClearScreen

			}

		case "ctrl+f":
			m.textinput.SetValue("")
			return model{
				table:       m.table,
				viewport:    m.viewport,
				textinput:   m.textinput,
				currentView: SearchContainer,
			}, tea.ClearScreen

		case "ctrl+l":
			containerLogs, err := dockerClient.ContainerLogs(m.table.SelectedRow()[0])
			if err != nil {
				panic(err)
			}

			headerHeight := lipgloss.Height(models.HeaderView(m.logsView.pager, m.table.SelectedRow()[1]))
			p := models.NewPager(widthScreen, heightScreen, containerLogs, headerHeight)

			lv := LogsView{
				pager:     p,
				container: m.table.SelectedRow()[1],
				image:     m.table.SelectedRow()[2],
			}

			return model{
				table:       m.table,
				viewport:    m.viewport,
				textinput:   m.textinput,
				logsView:    lv,
				currentView: LogsContainer,
			}, tea.ClearScreen

		case "ctrl+o":
			ov := models.NewOptions()
			return model{
				table:       m.table,
				viewport:    m.viewport,
				textinput:   m.textinput,
				logsView:    m.logsView,
				optionsView: ov,
				currentView: OptionsContainer,
			}, tea.ClearScreen

		case "down", "j":
			m.optionsView.Cursor++
			if m.optionsView.Cursor >= len(m.optionsView.Choices) {
				m.optionsView.Cursor = 0
			}

		case "up", "k":
			m.optionsView.Cursor--
			if m.optionsView.Cursor < 0 {
				m.optionsView.Cursor = len(m.optionsView.Choices) - 1
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(models.HeaderView(m.logsView.pager, ""))
		footerHeight := lipgloss.Height(models.FooterView(m.logsView.pager))
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			widthScreen = msg.Width
			heightScreen = msg.Height - verticalMarginHeight
			m.logsView.pager.YPosition = headerHeight
			m.ready = true
			m.logsView.pager.YPosition = headerHeight + 1
		} else {
			m.logsView.pager.Width = msg.Width
			m.logsView.pager.Height = msg.Height - verticalMarginHeight
		}

	}

	m.table, _ = m.table.Update(msg)
	m.viewport, _ = m.viewport.Update(msg)
	m.textinput, _ = m.textinput.Update(msg)
	m.logsView.pager, cmd = m.logsView.pager.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	switch m.currentView {
	case ListContainer:
		return baseStyle.Render(m.table.View()) + helpStyle("\n  ↑/↓: Navigate • Ctrl/C: Exit • Ctrl/F: Search • Ctrl/L: Logs container \n")
	case DetailContainer:
		return m.viewport.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	case SearchContainer:
		return fmt.Sprintf(
			"Search container by name\n\n%s\n\n%s",
			m.textinput.View(),
			"(esc to back)",
		) + "\n"
	case LogsContainer:
		return fmt.Sprintf("%s\n%s\n%s", models.HeaderView(m.logsView.pager, m.logsView.container+" - "+m.logsView.image), m.logsView.pager.View(), models.FooterView(m.logsView.pager))
	case OptionsContainer:
		return m.optionsView.View()
	default:
		return ""

	}

}

var dockerClient *docker.Docker
var containerList []docker.MyContainer
var widthScreen int
var heightScreen int

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
		table:       t,
		textinput:   models.NewTextInput(),
		currentView: ListContainer,
	}

	if _, err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
