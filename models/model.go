package models

import (
	"fmt"
	"os/exec"

	"github.com/ernesto27/dcli/docker"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const commands = `
 GENERAL ↑/↓: Navigate • ctrl+c: Exit • ctrl+r: refresh • esc: Back 
 CONTAINERS ctrl+f: Search • ctrl+l: Logs • ctrl+o: Options • ctrl+e: Attach cmd • ctrl+s: Stats
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
	MContainerStats
	MContainerExecOptions

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

	MStackList
	MStackDetail
)

type model struct {
	dockerClient         *docker.Docker
	containerList        ContainerList
	containerDetail      ContainerDetail
	containerSearch      ContainerSearch
	containerLogs        LogsView
	containerOptions     ContainerOptions
	containerStats       viewport.Model
	containerExecOptions ContainerExecOptions
	imageList            ImageList
	imageDetail          viewport.Model
	imageSearch          ImageSearch
	imageOptions         ImageOptions
	networkList          NetworkList
	networkSearch        NetworkSearch
	networkDetail        viewport.Model
	networkOptions       NetworkOptions
	volumeList           VolumeList
	volumeDetail         viewport.Model
	volumeSearch         VolumeSearch
	volumeOptions        VolumeOptions
	stackList            StackList
	stackDetail          viewport.Model
	ready                bool
	currentModel         currentModel
	ContainerID          string
	dockerVersion        string
	err                  error
	widthScreen          int
	heightScreen         int
	cpuCores             int
	ram                  string
}

func NewModel(dockerClient *docker.Docker, version string, cpuCores int, ram string) *model {
	m := &model{
		dockerClient:    dockerClient,
		containerSearch: NewContainerSearch(),
		imageSearch:     NewImageSearch(),
		networkSearch:   NewNetworkSearch(),
		volumeSearch:    NewVolumeSearch(),
		currentModel:    MContainerList,
		dockerVersion:   version,
		cpuCores:        cpuCores,
		ram:             ram,
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

			if m.currentModel == MStackDetail {
				m.currentModel = MStackList
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

		case "enter":
			if m.currentModel == MContainerExecOptions {
				m.currentModel = MContainerList
				switch m.containerExecOptions.Choices[m.containerExecOptions.Cursor] {
				case Bash:
					return m, attachToContainer(m.containerList.table.SelectedRow()[0], Bash)
				case Sh:
					return m, attachToContainer(m.containerList.table.SelectedRow()[0], Sh)
				case Ash:
					return m, attachToContainer(m.containerList.table.SelectedRow()[0], Ash)
				}
			}
		case "ctrl+p":
			stacks, err := m.dockerClient.StackList()
			if err != nil {
				fmt.Println(stacks)
			}
			m.stackList = NewStackList(stacks, "")
			m.currentModel = MStackList

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
	m.containerExecOptions, _ = m.containerExecOptions.Update(msg, &m)

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

	m.stackList.table, _ = m.stackList.Update(msg, &m)
	m.stackDetail, _ = m.stackDetail.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var statsStyle = lipgloss.NewStyle().
	MarginLeft(1).
	Background(lipgloss.Color("#021f2b")).
	Foreground(lipgloss.Color("#FAFAFA"))

var titleTableStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#740ceb")).
	Foreground(lipgloss.Color("#FFF")).
	MarginLeft(1).
	MarginTop(1).
	Padding(0, 1).
	Italic(true).
	Bold(true).
	Render

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
		return m.containerList.View(commands, &m)
	case MContainerDetail:
		return m.containerDetail.View()
	case MContainerSearch:
		return m.containerSearch.View()
	case MContainerLogs:
		return fmt.Sprintf("%s\n%s\n%s", HeaderView(m.containerLogs.pager, m.containerLogs.container+" - "+m.containerLogs.image), m.containerLogs.pager.View(), FooterView(m.containerLogs.pager))
	case MContainerOptions:
		return m.containerOptions.View()
	case MContainerStats:
		return m.containerStats.View()
	case MContainerExecOptions:
		return m.containerExecOptions.View()

	case MImageList:
		return m.imageList.View(commands, &m)
	case MImageOptions:
		return m.imageOptions.View()
	case MImageDetail:
		return m.imageDetail.View()
	case MImageSearch:
		return m.imageSearch.View()

	case MNetworkList:
		return m.networkList.View(commands, &m)
	case MNetworkSearch:
		return m.networkSearch.View()
	case MNetworkDetail:
		return m.networkDetail.View()
	case MNetworkOptions:
		return m.networkOptions.View()

	case MVolumeList:
		return m.volumeList.View(commands, &m)
	case MVolumeDetail:
		return m.volumeDetail.View()
	case MVolumeSearch:
		return m.volumeSearch.View()
	case MVolumeOptions:
		return m.volumeOptions.View()

	case MStackList:
		return m.stackList.View(commands, &m)
	case MStackDetail:
		return m.stackDetail.View()

	default:
		return ""

	}
}

type attachExited struct{ err error }

func attachToContainer(ID string, command string) tea.Cmd {
	c := exec.Command("docker", "exec", "-it", ID, command)
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

func (m *model) getDockerStats() string {
	return fmt.Sprintf("\U0001F433 DockerVersion: %s | Containers: %d (%s)| Images: %d (%s) | Volumes: %d  \U0001F5A5  CPU: %d | Memory: %s ",
		m.dockerVersion,
		len(m.dockerClient.Containers),
		m.dockerClient.GetAllContainersSize(),
		len(m.dockerClient.Images),
		m.dockerClient.GetAllImagesSize(),
		len(m.dockerClient.Volumes),
		m.cpuCores,
		m.ram,
	)
}

func (m *model) renderTable(title string, table string, commands string) string {
	dockerStats := statsStyle.Render(m.getDockerStats())
	return titleTableStyle(title) + "\n" + baseStyle.Render(table) + helpStyle("\n"+dockerStats+"\n"+commands)
}
