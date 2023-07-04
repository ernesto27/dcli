package models

import (
	"dockerniceui/docker"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const commands = `
 GENERAL ↑/↓: Navigate • ctrl+c: Exit • ctrl+r: refresh • esc: Back 
 CONTAINERS ctrl+f: Search • ctrl+l: Logs • ctrl+o: Options • ctrl+e: Attach cmd
 IMAGES ctrl+b: List • ctrl+f: Search • ctrl+o: Options
 NETWORKS ctrl+n: List • ctrl+f: Search  • ctrl+o: Options
 VOLUMES ctrl+v: List • ctrl+f: Search  • ctrl+o: Options
   `

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9999FF")).Render

type currentModel int

const (
	MContainerList currentModel = iota
	MContainerDetail
	MContainerSearch
	MContainerLogs
	MContainerOptions

	MImageList
	MImageDetail
	MImageSearch
	MImageOptions

	MNetworkList
	MNetworkSearch
	MNetworkDetail
	MNetworkOptions

	MVolumeList
	MVolumeDetail
	MVolumeSearch
	MVolumeOptions
)

type model struct {
	dockerClient     *docker.Docker
	containerList    ContainerList
	containerDetail  ContainerDetail
	containerSearch  ContainerSearch
	containerLogs    LogsView
	containerOptions ContainerOptions
	imageList        ImageList
	imageDetail      viewport.Model
	imageSearch      ImageSearch
	imageOptions     ImageOptions
	networkList      NetworkList
	networkSearch    NetworkSearch
	networkDetail    viewport.Model
	networkOptions   NetworkOptions
	volumeList       VolumeList
	volumeDetail     viewport.Model
	volumeSearch     VolumeSearch
	volumeOptions    VolumeOptions
	ready            bool
	currentModel     currentModel
	ContainerID      string
	dockerVersion    string
	err              error
	widthScreen      int
	heightScreen     int
}

func NewModel(dockerClient *docker.Docker, version string) *model {
	m := &model{
		dockerClient:    dockerClient,
		containerSearch: NewContainerSearch(),
		imageSearch:     NewImageSearch(),
		networkSearch:   NewNetworkSearch(),
		volumeSearch:    NewVolumeSearch(),
		currentModel:    MContainerList,
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
			if m.currentModel == MImageDetail || m.currentModel == MImageOptions {
				m.currentModel = MImageList
				return m, tea.ClearScreen
			}

			if m.currentModel == MNetworkDetail || m.currentModel == MNetworkSearch || m.currentModel == MNetworkOptions {
				m.currentModel = MNetworkList
				return m, tea.ClearScreen
			}

			if m.currentModel == MVolumeDetail || m.currentModel == MVolumeSearch || m.currentModel == MVolumeOptions {
				m.currentModel = MVolumeList
				return m, tea.ClearScreen
			}

		case "ctrl+c":
			return m, tea.Quit
		case "down":
			m.containerOptions.Cursor++
			if m.containerOptions.Cursor >= len(m.containerOptions.Choices) {
				m.containerOptions.Cursor = 0
			}

		case "up":
			m.containerOptions.Cursor--
			if m.containerOptions.Cursor < 0 {
				m.containerOptions.Cursor = len(m.containerOptions.Choices) - 1
			}

		case "ctrl+e":
			return m, attachToContainer(m.containerList.table.SelectedRow()[0])

		case "ctrl+r":
			m.setContainerList()
			return m, tea.ClearScreen

		case "ctrl+v":
			volumes, err := m.dockerClient.VolumeList()
			if err != nil {
				fmt.Println(err)
			}
			m.volumeList = NewVolumeList(volumes, "")
			m.currentModel = MVolumeList

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

	m.containerList.table, _ = m.containerList.Update(msg, &m)
	m.containerDetail.viewport, _ = m.containerDetail.Update(msg, &m)
	m.containerSearch, _ = m.containerSearch.Update(msg, &m)
	m.containerOptions, _ = m.containerOptions.Update(msg, &m)
	m.containerLogs.pager, _ = m.containerLogs.pager.Update(msg)

	m.imageList.table, _ = m.imageList.Update(msg, &m)
	m.imageSearch, _ = m.imageSearch.Update(msg, &m)
	m.imageOptions, _ = m.imageOptions.Update(msg, &m)
	m.imageDetail, _ = m.imageDetail.Update(msg)

	m.networkList.table, _ = m.networkList.Update(msg, &m)
	m.networkSearch, _ = m.networkSearch.Update(msg, &m)
	m.networkOptions, _ = m.networkOptions.Update(msg, &m)

	m.volumeList.table, _ = m.volumeList.Update(msg, &m)
	m.volumeDetail, _ = m.volumeDetail.Update(msg)
	m.volumeSearch, _ = m.volumeSearch.Update(msg, &m)
	m.volumeOptions, _ = m.volumeOptions.Update(msg, &m)

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

	switch m.currentModel {
	case MContainerList:
		return m.containerList.View(commands, m.dockerVersion)
	case MContainerDetail:
		return m.containerDetail.View()
	case MContainerSearch:
		return m.containerSearch.View()
	case MContainerLogs:
		return fmt.Sprintf("%s\n%s\n%s", HeaderView(m.containerLogs.pager, m.containerLogs.container+" - "+m.containerLogs.image), m.containerLogs.pager.View(), FooterView(m.containerLogs.pager))
	case MContainerOptions:
		return m.containerOptions.View()

	case MImageList:
		return m.imageList.View(commands, m.dockerVersion)
	case MImageOptions:
		return m.imageOptions.View()
	case MImageDetail:
		return m.imageDetail.View()
	case MImageSearch:
		return m.imageSearch.View()

	case MNetworkList:
		return m.networkList.View(commands, m.dockerVersion)
	case MNetworkSearch:
		return m.networkSearch.View()
	case MNetworkDetail:
		return m.networkDetail.View()
	case MNetworkOptions:
		return m.networkOptions.View()

	case MVolumeList:
		return m.volumeList.View(commands, m.dockerVersion)
	case MVolumeDetail:
		return m.volumeDetail.View()
	case MVolumeSearch:
		return m.volumeSearch.View()
	case MVolumeOptions:
		return m.volumeOptions.View()

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
	m.currentModel = MContainerList
}
