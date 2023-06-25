package main

import (
	"context"
	"dockerniceui/docker"
	"dockerniceui/models"
	"fmt"
	"os"

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

	m := models.NewModel(dockerClient, version)
	if _, err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
