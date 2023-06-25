package main

import (
	"context"
	"dockerniceui/docker"
	"dockerniceui/models"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

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

	dockerClient.Events()

	// m := model{
	// 	containerList:   getTableWithData(),
	// 	containerSearch: models.NewSearch(),
	// 	imageSearch:     models.NewSearch(),
	// 	currentView:     ContainerList,
	// 	dockerVersion:   version,
	// }

	m := models.NewModel(&dockerClient, version)

	if _, err := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
