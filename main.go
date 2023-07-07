package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/ernesto27/dcli/docker"
	"github.com/ernesto27/dcli/models"
	"github.com/ernesto27/dcli/utils"
	"github.com/shirou/gopsutil/mem"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := context.Background()
	dockerClient, err := docker.New(ctx)
	if err != nil {
		panic(err)
	}

	version, err := dockerClient.ServerVersion()
	if err != nil {
		fmt.Println(err)
	}

	dockerClient.Events()
	// set images, volumes for count view
	dockerClient.ImageList()
	dockerClient.VolumeList()

	ram := ""
	v, err := mem.VirtualMemory()
	if err == nil {
		ram = utils.FormatSize(int64(v.Total))
	}

	m := models.NewModel(dockerClient, version, runtime.NumCPU(), ram)
	if _, err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
