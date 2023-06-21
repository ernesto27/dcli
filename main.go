package main

import (
	"context"
	"dockerniceui/docker"
	"dockerniceui/models"
	"dockerniceui/utils"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9999FF")).Render

type currentView int

const (
	ContainerList currentView = iota
	ContainerDetail
	ContainerSearch
	ContainerLogs
	ContainerOptions
	ImageList
	ImageDetail
)

type LogsView struct {
	pager     viewport.Model
	container string
	image     string
}

type model struct {
	table         table.Model
	viewport      viewport.Model
	textinput     textinput.Model
	logsView      LogsView
	optionsView   models.Options
	imageTable    table.Model
	imageDetail   viewport.Model
	ready         bool
	currentView   currentView
	ContainerID   string
	dockerVersion string
	err           error
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
			m.err = nil
			t := models.NewTable(models.GetContainerColumns(), models.GetContainerRows(dockerClient.Containers, ""))
			m.table = t
			m.currentView = ContainerList
			return m, tea.ClearScreen

		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.currentView {
			case ContainerList:
				container, err := dockerClient.GetContainerByName(m.table.SelectedRow()[1])
				if err != nil {
					fmt.Println(err)
				}

				vp, err := models.NewViewport(utils.GetContent(container, utils.CreateTable))
				if err != nil {
					fmt.Println(err)
				}

				m.viewport = vp
				m.currentView = ContainerDetail
				return m, tea.ClearScreen
			case ContainerSearch:
				value := m.textinput.Value()
				t := models.NewTable(models.GetContainerColumns(), models.GetContainerRows(dockerClient.Containers, value))
				m.table = t
				m.currentView = ContainerList
				return m, tea.ClearScreen
			case ContainerOptions:
				errAction := false
				switch m.optionsView.Choices[m.optionsView.Cursor] {
				case models.Stop:
					err := dockerClient.ContainerStop(m.ContainerID)
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
					time.Sleep(1 * time.Second)
				case models.Start:
					err := dockerClient.ContainerStart(m.ContainerID)
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
				case models.Remove:
					err := dockerClient.ContainerRemove(m.ContainerID)
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
				}

				if !errAction {
					t := getTableWithData()
					m.table = t
					m.currentView = ContainerList
					return m, tea.ClearScreen
				}

			case ImageList:
				img, err := dockerClient.GetImageByID(m.imageTable.SelectedRow()[0])
				if err != nil {
					fmt.Println(err)
				}

				imgView, err := models.NewImageDetail(utils.GetContentDetailImage(img, utils.CreateTable))
				if err != nil {
					fmt.Println(err)
				}

				m.imageDetail = imgView
				m.currentView = ImageDetail

				return m, tea.ClearScreen
			}

		case "ctrl+f":
			m.textinput.SetValue("")
			m.currentView = ContainerSearch
			return m, tea.ClearScreen

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

			m.logsView = lv
			m.currentView = ContainerLogs

			return m, tea.ClearScreen

		case "ctrl+o":
			ov := models.NewOptions(m.table.SelectedRow()[1], m.table.SelectedRow()[2])
			m.optionsView = ov
			m.currentView = ContainerOptions
			m.ContainerID = m.table.SelectedRow()[0]
			return m, tea.ClearScreen

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

		case "ctrl+e":
			return m, attachToContainer(m.table.SelectedRow()[0])

		case "ctrl+b":
			images, err := dockerClient.ImageList()
			if err != nil {
				fmt.Println(err)
			}

			m.imageTable = models.NewImageTable(models.GetImageRows(images, ""))
			m.currentView = ImageList
			return m, tea.ClearScreen
		}

	case attachExited:
		if msg.err != nil {
			m.err = msg.err
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
	m.imageTable, _ = m.imageTable.Update(msg)
	m.logsView.pager, cmd = m.logsView.pager.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#F1ECEB")).
			Foreground(lipgloss.Color("#FC765B"))

		errorText := ""
		if m.err.Error() == "exit status 1" {
			errorText = "Container is not running"
		} else {
			errorText = "OCI runtime exec failed: exec failed: unable to start container process: exec"
		}

		return errorStyle.Render("Error: " + errorText + " \n\nEsc to go back")
	}

	commands := `
 GENERAL ↑/↓: Navigate • Ctrl/C: Exit • Esc: Back 
 CONTAINERS Ctrl/F: Search • Ctrl/L: Logs • Ctrl/O: Options • Ctrl/E: Attach cmd
 IMAGES Ctrl/B: List
	`

	switch m.currentView {
	case ContainerList:
		return baseStyle.Render(m.table.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
	case ContainerDetail:
		return m.viewport.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	case ContainerSearch:
		return fmt.Sprintf(
			"Search container by name\n\n%s\n\n%s",
			m.textinput.View(),
			"(esc to back)",
		) + "\n"
	case ContainerLogs:
		return fmt.Sprintf("%s\n%s\n%s", models.HeaderView(m.logsView.pager, m.logsView.container+" - "+m.logsView.image), m.logsView.pager.View(), models.FooterView(m.logsView.pager))
	case ContainerOptions:
		return m.optionsView.View()

	case ImageList:
		return baseStyle.Render(m.imageTable.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
	case ImageDetail:
		return m.imageDetail.View()
	default:
		return ""

	}

}

type attachExited struct{ err error }

func attachToContainer(ID string) tea.Cmd {
	c := exec.Command("docker", "exec", "-it", ID, "bin/bash")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return attachExited{err}
	})
}

var dockerClient *docker.Docker
var widthScreen int
var heightScreen int

func main() {
	ctx := context.Background()
	var err error
	dockerClient, err = docker.New(ctx)
	if err != nil {
		panic(err)
	}

	version, err := dockerClient.ServerVersion()
	if err != nil {
		fmt.Println(err)
	}

	m := model{
		table:         getTableWithData(),
		textinput:     models.NewTextInput(),
		currentView:   ContainerList,
		dockerVersion: version,
	}

	if _, err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getTableWithData() table.Model {
	var err error
	_, err = dockerClient.ContainerList()
	// TODO FIX
	if err != nil {
		panic(err)
	}

	t := models.NewTable(models.GetContainerColumns(), models.GetContainerRows(dockerClient.Containers, ""))
	return t
}
