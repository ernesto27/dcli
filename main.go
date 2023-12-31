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

	msgError := "Error connecting to docker daemon, please check if docker is installed and running..."
	dockerClient, err := docker.New(ctx)
	if err != nil {
		fmt.Println(msgError)
		os.Exit(1)
	}

	version, err := dockerClient.ServerVersion()
	if err != nil {
		fmt.Println(msgError)
		os.Exit(1)
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
