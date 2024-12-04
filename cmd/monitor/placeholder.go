package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PlaceholderModel struct {
	Text string
}

func (p *PlaceholderModel) Refresh() {}

func (p *PlaceholderModel) Update(tea.Msg) tea.Cmd { return nil }

func (p *PlaceholderModel) View() string {
	style := lipgloss.NewStyle().Height(12)
	return style.Render(p.Text)
}
