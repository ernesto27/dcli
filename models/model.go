package models

import (
	"dockerniceui/docker"
	"dockerniceui/utils"
	"fmt"
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
	ImageOptions

	NetworkList
	NetworkSearch
)

type LogsView struct {
	pager     viewport.Model
	container string
	image     string
}

type model struct {
	dockerClient     *docker.Docker
	containerList    table.Model
	containerDetail  viewport.Model
	containerSearch  textinput.Model
	containerLogs    LogsView
	containerOptions Options
	imageList        table.Model
	imageDetail      viewport.Model
	imageSearch      textinput.Model
	imageOptions     Options
	networkList      table.Model
	networkSearch    textinput.Model
	ready            bool
	currentView      currentView
	ContainerID      string
	dockerVersion    string
	err              error
	widthScreen      int
	heightScreen     int
}

func NewModel(dockerClient *docker.Docker, version string) *model {
	m := &model{
		dockerClient:    dockerClient,
		containerSearch: NewSearch(),
		imageSearch:     NewSearch(),
		networkSearch:   NewSearch(),
		currentView:     ContainerList,
		dockerVersion:   version,
	}
	m.setContainerList()

	return m
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
			if m.currentView == ImageDetail || m.currentView == ImageOptions {
				m.currentView = ImageList
				return m, tea.ClearScreen
			}

			m.setContainerList()
			return m, tea.ClearScreen

		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.currentView {
			case ContainerList:
				container, err := m.dockerClient.GetContainerByName(m.containerList.SelectedRow()[1])
				if err != nil {
					fmt.Println(err)
				}

				vp, err := NewContainerDetail(container, utils.CreateTable)
				if err != nil {
					fmt.Println(err)
				}

				m.containerDetail = vp
				m.currentView = ContainerDetail
				return m, tea.ClearScreen
			case ContainerSearch:
				value := m.containerSearch.Value()
				t := NewContainerList(GetContainerRows(m.dockerClient.Containers, value))
				m.containerList = t
				m.currentView = ContainerList
				return m, tea.ClearScreen
			case ContainerOptions:
				errAction := false
				switch m.containerOptions.Choices[m.containerOptions.Cursor] {
				case Stop:
					err := m.dockerClient.ContainerStop(m.ContainerID)
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
					time.Sleep(1 * time.Second)
				case Start:
					err := m.dockerClient.ContainerStart(m.ContainerID)
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
				case Remove:
					err := m.dockerClient.ContainerRemove(m.ContainerID)
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
				case Restart:
					err := m.dockerClient.ContainerRestart(m.ContainerID)
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
				}

				if !errAction {
					m.setContainerList()
					m.currentView = ContainerList
					return m, tea.ClearScreen
				}

			case ImageList:
				img, err := m.dockerClient.GetImageByID(m.imageList.SelectedRow()[0])
				if err != nil {
					fmt.Println(err)
				}

				imgView, err := NewImageDetail(img, utils.CreateTable)
				if err != nil {
					fmt.Println(err)
				}

				m.imageDetail = imgView
				m.currentView = ImageDetail

				return m, tea.ClearScreen

			case ImageSearch:
				value := m.imageSearch.Value()
				imgList := NewImageList(GetImageRows(m.dockerClient.Images, value))
				m.imageList = imgList
				m.currentView = ImageList
				return m, tea.ClearScreen

			case ImageOptions:
				errAction := false
				option := m.imageOptions.Choices[m.imageOptions.Cursor]

				if option == Remove {
					err := m.dockerClient.ImageRemove(m.imageList.SelectedRow()[1])
					if err != nil {
						fmt.Println(err)
						errAction = true
					}
				}

				if !errAction {
					images, err := m.dockerClient.ImageList()
					if err != nil {
						fmt.Println(err)
					}

					m.imageList = NewImageList(GetImageRows(images, ""))
					m.currentView = ImageList
					return m, tea.ClearScreen
				}
			case NetworkSearch:
				value := m.networkSearch.Value()
				networks := NewNetworkList(GetNetworkRows(m.dockerClient.Networks, value))
				m.networkList = networks
				m.currentView = NetworkList
			}

		case "ctrl+f":
			switch m.currentView {
			case ContainerList:
				m.containerSearch.SetValue("")
				m.currentView = ContainerSearch
			case ImageList:
				m.imageSearch.SetValue("")
				m.currentView = ImageSearch
			case NetworkList:
				m.networkSearch.SetValue("")
				m.currentView = NetworkSearch
			}

			return m, tea.ClearScreen

		case "ctrl+l":
			if m.currentView == ContainerList {
				containerLogs, err := m.dockerClient.ContainerLogs(m.containerList.SelectedRow()[0])
				if err != nil {
					panic(err)
				}

				headerHeight := lipgloss.Height(HeaderView(m.containerLogs.pager, m.containerList.SelectedRow()[1]))
				p := NewContainerLogs(m.widthScreen, m.heightScreen, containerLogs, headerHeight)

				lv := LogsView{
					pager:     p,
					container: m.containerList.SelectedRow()[1],
					image:     m.containerList.SelectedRow()[2],
				}

				m.containerLogs = lv
				m.currentView = ContainerLogs

				return m, tea.ClearScreen
			}

		case "ctrl+o":
			if m.currentView == ContainerList {
				ov := NewContainerOptions(m.containerList.SelectedRow()[1], m.containerList.SelectedRow()[2])
				m.containerOptions = ov
				m.currentView = ContainerOptions
				m.ContainerID = m.containerList.SelectedRow()[0]
				return m, tea.ClearScreen
			} else if m.currentView == ImageList {
				ov := NewContainerOptions("", m.imageList.SelectedRow()[1])
				m.imageOptions = ov
				m.currentView = ImageOptions
				return m, tea.ClearScreen
			}

		case "down", "j":
			m.containerOptions.Cursor++
			if m.containerOptions.Cursor >= len(m.containerOptions.Choices) {
				m.containerOptions.Cursor = 0
			}

		case "up", "k":
			m.containerOptions.Cursor--
			if m.containerOptions.Cursor < 0 {
				m.containerOptions.Cursor = len(m.containerOptions.Choices) - 1
			}

		case "ctrl+e":
			return m, attachToContainer(m.containerList.SelectedRow()[0])

		case "ctrl+b":
			images, err := m.dockerClient.ImageList()
			if err != nil {
				fmt.Println(err)
			}

			m.imageList = NewImageList(GetImageRows(images, ""))
			m.currentView = ImageList
			return m, tea.ClearScreen

		case "ctrl+n":
			networks, err := m.dockerClient.NetworkList()
			if err != nil {
				fmt.Println(err)
			}

			m.networkList = NewNetworkList(GetNetworkRows(networks, ""))
			m.currentView = NetworkList
			return m, tea.ClearScreen

		case "ctrl+r":
			m.setContainerList()
			return m, tea.ClearScreen

		}

	case attachExited:
		if msg.err != nil {
			m.err = msg.err
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(HeaderView(m.containerLogs.pager, ""))
		footerHeight := lipgloss.Height(FooterView(m.containerLogs.pager))
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.widthScreen = msg.Width
			m.heightScreen = msg.Height - verticalMarginHeight
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
	m.imageDetail, _ = m.imageDetail.Update(msg)
	m.networkList, _ = m.networkList.Update(msg)
	m.networkSearch, _ = m.networkSearch.Update(msg)
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
 GENERAL ↑/↓: Navigate • ctrl/c: Exit • ctrl/r: refresh • esc: Back 
 CONTAINERS ctrl/f: Search • ctrl/l: Logs • ctrl/o: Options • ctrl/e: Attach cmd
 IMAGES ctrl/b: List • ctrl/f: Search
 NETWORKS ctrl/n: List
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
		return fmt.Sprintf("%s\n%s\n%s", HeaderView(m.containerLogs.pager, m.containerLogs.container+" - "+m.containerLogs.image), m.containerLogs.pager.View(), FooterView(m.containerLogs.pager))
	case ContainerOptions:
		return m.containerOptions.View()
	case ImageOptions:
		return m.imageOptions.View()

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

	case NetworkList:
		return baseStyle.Render(m.networkList.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
	case NetworkSearch:
		return fmt.Sprintf(
			"Search network by name\n\n%s\n\n%s",
			m.networkSearch.View(),
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

func (m *model) setContainerList() {
	var err error
	_, err = m.dockerClient.ContainerList()
	if err != nil {
		fmt.Println(err)
	}

	t := NewContainerList(GetContainerRows(m.dockerClient.Containers, ""))
	m.err = nil
	m.containerList = t
	m.currentView = ContainerList
}
