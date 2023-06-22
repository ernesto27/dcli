package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

func NewPager(width int, height int, logs string, headerHeight int) viewport.Model {
	p := viewport.New(width, height)
	p.YPosition = headerHeight + 1
	p.SetContent(logs)
	return p
}

func HeaderView(pager viewport.Model, text string) string {
	title := titleStyle.Render(text)
	line := strings.Repeat("─", max(0, pager.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func FooterView(pager viewport.Model) string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", pager.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, pager.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
