package models

import (
	"dockerniceui/docker"
	"dockerniceui/utils"
	"fmt"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
)

type LogsView struct {
	pager     viewport.Model
	container string
	image     string
}

type model struct {
	dockerClient     *docker.Docker
	containerList    ContainerList
	containerDetail  ContianerDetail
	containerSearch  Search
	containerLogs    LogsView
	containerOptions Options
	imageList        table.Model
	imageDetail      viewport.Model
	imageSearch      Search
	imageOptions     Options
	networkList      table.Model
	networkSearch    Search
	networkDetail    viewport.Model
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
		containerSearch: NewSearch(),
		imageSearch:     NewSearch(),
		networkSearch:   NewSearch(),
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

		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.currentModel {
			case MContainerOptions:
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
					m.currentModel = MContainerList
					return m, tea.ClearScreen
				}

			case MImageList:
				img, err := m.dockerClient.GetImageByID(m.imageList.SelectedRow()[0])
				if err != nil {
					fmt.Println(err)
				}

				imgView, err := NewImageDetail(img, utils.CreateTable)
				if err != nil {
					fmt.Println(err)
				}

				m.imageDetail = imgView
				m.currentModel = MImageDetail

				return m, tea.ClearScreen

			case MImageSearch:
				value := m.imageSearch.textInput.Value()
				imgList := NewImageList(GetImageRows(m.dockerClient.Images, value))
				m.imageList = imgList
				m.currentModel = MImageList
				return m, tea.ClearScreen

			case MImageOptions:
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
					m.currentModel = MImageList
					return m, tea.ClearScreen
				}
			case MNetworkSearch:
				value := m.networkSearch.textInput.Value()
				networks := NewNetworkList(GetNetworkRows(m.dockerClient.Networks, value))
				m.networkList = networks
				m.currentModel = MNetworkList
			case MNetworkList:
				network, err := m.dockerClient.GetNetworkByName(m.networkList.SelectedRow()[1])
				if err != nil {
					fmt.Println(err)
				}

				nd, err := NewNetworkDetail(network, utils.CreateTable)
				if err != nil {
					fmt.Println(err)
				}
				m.networkDetail = nd
				m.currentModel = MNetworkDetail

			}

		case "ctrl+f":
			switch m.currentModel {
			case MContainerList:
				m.containerSearch.textInput.SetValue("")
				m.currentModel = MContainerSearch
			case MImageList:
				m.imageSearch.textInput.SetValue("")
				m.currentModel = MImageSearch
			case MNetworkList:
				m.networkSearch.textInput.SetValue("")
				m.currentModel = MNetworkSearch
			}

			return m, tea.ClearScreen

		case "ctrl+l":
			if m.currentModel == MContainerList {
				containerLogs, err := m.dockerClient.ContainerLogs(m.containerList.table.SelectedRow()[0])
				if err != nil {
					panic(err)
				}

				headerHeight := lipgloss.Height(HeaderView(m.containerLogs.pager, m.containerList.table.SelectedRow()[1]))
				p := NewContainerLogs(m.widthScreen, m.heightScreen, containerLogs, headerHeight)

				lv := LogsView{
					pager:     p,
					container: m.containerList.table.SelectedRow()[1],
					image:     m.containerList.table.SelectedRow()[2],
				}

				m.containerLogs = lv
				m.currentModel = MContainerLogs

				return m, tea.ClearScreen
			}

		case "ctrl+o":
			if m.currentModel == MContainerList {
				ov := NewContainerOptions(m.containerList.table.SelectedRow()[1], m.containerList.table.SelectedRow()[2])
				m.containerOptions = ov
				m.currentModel = MContainerOptions
				m.ContainerID = m.containerList.table.SelectedRow()[0]
				return m, tea.ClearScreen
			} else if m.currentModel == MImageList {
				ov := NewContainerOptions("", m.imageList.SelectedRow()[1])
				m.imageOptions = ov
				m.currentModel = MImageOptions
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
			return m, attachToContainer(m.containerList.table.SelectedRow()[0])

		case "ctrl+b":
			images, err := m.dockerClient.ImageList()
			if err != nil {
				fmt.Println(err)
			}

			m.imageList = NewImageList(GetImageRows(images, ""))
			m.currentModel = MImageList
			return m, tea.ClearScreen

		case "ctrl+n":
			networks, err := m.dockerClient.NetworkList()
			if err != nil {
				fmt.Println(err)
			}

			m.networkList = NewNetworkList(GetNetworkRows(networks, ""))
			m.currentModel = MNetworkList
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

	m.containerList.table, _ = m.containerList.Update(msg, &m)
	m.containerDetail.viewport, _ = m.containerDetail.Update(msg, &m)
	m.containerSearch.textInput, _ = m.containerSearch.Update(msg, &m)
	m.containerLogs.pager, _ = m.containerLogs.pager.Update(msg)
	m.imageList, _ = m.imageList.Update(msg)
	m.imageDetail, _ = m.imageDetail.Update(msg)
	m.networkList, _ = m.networkList.Update(msg)
	m.networkSearch.textInput, _ = m.networkSearch.textInput.Update(msg)
	m.imageSearch.textInput, cmd = m.imageSearch.textInput.Update(msg)

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

	switch m.currentModel {
	case MContainerList:
		// return baseStyle.Render(m.containerList.table.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
		return m.containerList.View()
	case MContainerDetail:
		return m.containerDetail.View()
		// return m.containerDetail.View() + helpStyle("\n  ↑/↓: Navigate • Esc: back to list\n")
	case MContainerSearch:
		return m.containerSearch.View()
		// return fmt.Sprintf(
		// 	"Search container by name\n\n%s\n\n%s",
		// 	m.containerSearch.textInput.View(),
		// 	"(esc to back)",
		// ) + "\n"
	case MContainerLogs:
		return fmt.Sprintf("%s\n%s\n%s", HeaderView(m.containerLogs.pager, m.containerLogs.container+" - "+m.containerLogs.image), m.containerLogs.pager.View(), FooterView(m.containerLogs.pager))
	case MContainerOptions:
		return m.containerOptions.View()
	case MImageOptions:
		return m.imageOptions.View()

	case MImageList:
		return baseStyle.Render(m.imageList.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
	case MImageDetail:
		return m.imageDetail.View()
	case MImageSearch:
		return fmt.Sprintf(
			"Search image by name\n\n%s\n\n%s",
			m.imageSearch.textInput.View(),
			"(esc to back)",
		) + "\n"

	case MNetworkList:
		return baseStyle.Render(m.networkList.View()) + helpStyle("\n DockerVersion: "+m.dockerVersion+" \n"+commands)
	case MNetworkSearch:
		return fmt.Sprintf(
			"Search network by name\n\n%s\n\n%s",
			m.networkSearch.textInput.View(),
			"(esc to back)",
		) + "\n"
	case MNetworkDetail:
		return m.networkDetail.View()

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
