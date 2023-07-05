package models

import (
	"dcli/docker"
	"dcli/utils"
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func NewContainerStats(stats docker.MyContainerStats, createTable utils.CreateTableFunc) (viewport.Model, error) {
	content := getContentStats(stats)
	const width = 120

	vp := viewport.New(width, 30)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return viewport.Model{}, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return viewport.Model{}, err
	}

	vp.SetContent(str)

	return vp, nil
}

func getContentStats(stats docker.MyContainerStats) string {
	response := ""

	response += utils.CreateTable("# Stats", []string{"CPU", "MEM USAGE/LIMIT", "MEM", "PIDS"},
		[][]string{
			{
				fmt.Sprintf("%.2f%%", stats.CPUPer),
				fmt.Sprintf("%s / %s", stats.MemUsage, stats.MemLimit),
				fmt.Sprintf("%.2f%%", stats.MemPer),
				fmt.Sprintf("%d", stats.PID),
			},
		})

	return response
}
