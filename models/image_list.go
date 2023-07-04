package models

import (
	"dockerniceui/docker"
	"dockerniceui/utils"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ImageList struct {
	table table.Model
}

func NewImageList(images []docker.MyImage, query string) ImageList {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Image", Width: 40},
		{Title: "Size", Width: 20},
		{Title: "Created", Width: 20},
	}

	rows := GetImageRows(images, query)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(180),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return ImageList{table: t}
}

func (il ImageList) View(commands string, dockerVersion string) string {
	return baseStyle.Render(il.table.View()) + helpStyle("\n DockerVersion: "+dockerVersion+" \n"+commands)
}

func (il ImageList) Update(msg tea.Msg, m *model) (table.Model, tea.Cmd) {
	il.table, _ = il.table.Update(msg)
	if m.currentModel != MImageList {
		return il.table, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			img, err := m.dockerClient.GetImageByID(m.imageList.table.SelectedRow()[0])
			if err != nil {
				fmt.Println(err)
			}

			imgView, err := NewImageDetail(img, utils.CreateTable)
			if err != nil {
				fmt.Println(err)
			}

			m.imageDetail = imgView
			m.currentModel = MImageDetail
		case "ctrl+f":
			m.imageSearch.textInput.SetValue("")
			m.currentModel = MImageSearch
		case "ctrl+o":
			ov := NewImageOptions(m.imageList.table.SelectedRow()[1])
			m.imageOptions = ov
			m.currentModel = MImageOptions

		}
	}

	return il.table, nil

}

func GetImageRows(images []docker.MyImage, query string) []table.Row {
	var filtered []docker.MyImage
	if query == "" {
		filtered = images
	} else {
		for _, i := range images {
			if strings.Contains(strings.ToLower(i.Summary.RepoTags[0]), strings.ToLower(query)) {
				filtered = append(filtered, i)
			}
		}
	}

	rowsItems := []table.Row{}
	for _, i := range filtered {
		item := []string{i.Summary.ID, i.Summary.RepoTags[0], i.GetFormatSize(), i.GetFormatTimestamp()}
		rowsItems = append(rowsItems, item)
	}

	return rowsItems
}
