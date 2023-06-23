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
	ImageSearch
)

type LogsView struct {
	pager     viewport.Model
	container string
	image     string
}

type model struct {
	containerList   table.Model
	containerDetail viewport.Model
	containerSearch textinput.Model
	containerLogs   LogsView
	optionsView     models.Options
	imageList       table.Model
	imageDetail     viewport.Model
	imageSearch     textinput.Model
	ready           bool
	currentView     currentView
	ContainerID     string
	dockerVersion   string
	err             error
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
			if m.currentView == ImageDetail {
				m.currentView = ImageList
				return m, tea.ClearScreen
			}

			m.err = nil
			t := models.NewContainerList(models.GetContainerRows(dockerClient.Containers, ""))
			m.containerList = t
			m.currentView = ContainerList
			return m, tea.ClearScreen

		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.currentView {
			case ContainerList:
				container, err := dockerClient.GetContainerByName(m.containerList.SelectedRow()[1])
				if err != nil {
					fmt.Println(err)
				}

				vp, err := models.NewContainerDetail(container, utils.CreateTable)
				if err != nil {
					fmt.Println(err)
				}

				m.containerDetail = vp
				m.currentView = ContainerDetail
				return m, tea.ClearScreen
			case ContainerSearch:
				value := m.containerSearch.Value()
				t := models.NewContainerList(models.GetContainerRows(dockerClient.Containers, value))
				m.containerList = t
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
					m.containerList = t
					m.currentView = ContainerList
					return m, tea.ClearScreen
				}

			case ImageList:
				img, err := dockerClient.GetImageByID(m.imageList.SelectedRow()[0])
				if err != nil {
					fmt.Println(err)
				}

				imgView, err := models.NewImageDetail(img, utils.CreateTable)
				if err != nil {
					fmt.Println(err)
				}

				m.imageDetail = imgView
				m.currentView = ImageDetail

				return m, tea.ClearScreen

			case ImageSearch:
				value := m.imageSearch.Value()
				imgList := models.NewImageList(models.GetImageRows(dockerClient.Images, value))
				m.imageList = imgList
				m.currentView = ImageList
				return m, tea.ClearScreen
			}

		case "ctrl+f":
			if m.currentView == ContainerList {
				m.containerSearch.SetValue("")
				m.currentView = ContainerSearch

			} else if m.currentView == ImageList {
				m.imageSearch.SetValue("")
				m.currentView = ImageSearch
			}

			return m, tea.ClearScreen

		case "ctrl+l":
			containerLogs, err := dockerClient.ContainerLogs(m.containerList.SelectedRow()[0])
			if err != nil {
				panic(err)
			}

			headerHeight := lipgloss.Height(models.HeaderView(m.containerLogs.pager, m.containerList.SelectedRow()[1]))
			p := models.NewContainerLogs(widthScreen, heightScreen, containerLogs, headerHeight)

			lv := LogsView{
				pager:     p,
				container: m.containerList.SelectedRow()[1],
				image:     m.containerList.SelectedRow()[2],
			}

			m.containerLogs = lv
			m.currentView = ContainerLogs

			return m, tea.ClearScreen

		case "ctrl+o":
			ov := models.NewContainerOptions(m.containerList.SelectedRow()[1], m.containerList.SelectedRow()[2])
			m.optionsView = ov
			m.currentView = ContainerOptions
			m.ContainerID = m.containerList.SelectedRow()[0]
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
			return m, attachToContainer(m.containerList.SelectedRow()[0])

		case "ctrl+b":
			images, err := dockerClient.ImageList()
			if err != nil {
				fmt.Println(err)
			}

			m.imageList = models.NewImageList(models.GetImageRows(images, ""))
			m.currentView = ImageList
			return m, tea.ClearScreen
		}

	case attachExited:
		if msg.err != nil {
			m.err = msg.err
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(models.HeaderView(m.containerLogs.pager, ""))
		footerHeight := lipgloss.Height(models.FooterView(m.containerLogs.pager))
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			widthScreen = msg.Width
			heightScreen = msg.Height - verticalMarginHeight
			m.containerLogs.pager.YPosition = headerHeight
			m.ready = true
			m.containerLogs.pager.YPosition = headerHeight + 1
		} else {
			m.containerLogs.pager.Width = msg.Width
			m.containerLogs.pager.Height = msg.Height - verticalMarginHeight
		}

	}

	m.containerList, _ = m.containerList.Update(msg)
	m.containerDetail, _ = m.containerDetail.Update(msg)
	m.containerSearch, _ = m.containerSearch.Update(msg)
	m.containerLogs.pager, _ = m.containerLogs.pager.Update(msg)
	m.imageList, _ = m.imageList.Update(msg)
	m.imageSearch, cmd = m.imageSearch.Update(msg)
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
 GENERAL ↑/↓: Navigate • ctrl/c: Exit • esc: Back 
 CONTAINERS ctrl/f: Search • ctrl/l: Logs • ctrl/o: Options • ctrl/e: Attach cmd
 IMAGES ctrl/b: List • ctrl/f: Search
	`

	switch m.currentView {
	case ContainerList:
		return baseStyle.Render(m.containerList.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
	case ContainerDetail:
		return m.containerDetail.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	case ContainerSearch:
		return fmt.Sprintf(
			"Search container by name\n\n%s\n\n%s",
			m.containerSearch.View(),
			"(esc to back)",
		) + "\n"
	case ContainerLogs:
		return fmt.Sprintf("%s\n%s\n%s", models.HeaderView(m.containerLogs.pager, m.containerLogs.container+" - "+m.containerLogs.image), m.containerLogs.pager.View(), models.FooterView(m.containerLogs.pager))
	case ContainerOptions:
		return m.optionsView.View()

	case ImageList:
		return baseStyle.Render(m.imageList.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
	case ImageDetail:
		return m.imageDetail.View()
	case ImageSearch:
		return fmt.Sprintf(
			"Search image by name\n\n%s\n\n%s",
			m.imageSearch.View(),
			"(esc to back)",
		) + "\n"
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
		containerList:   getTableWithData(),
		containerSearch: models.NewSearch(),
		imageSearch:     models.NewSearch(),
		currentView:     ContainerList,
		dockerVersion:   version,
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

	t := models.NewContainerList(models.GetContainerRows(dockerClient.Containers, ""))
	return t
}
