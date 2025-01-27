package main

import (
	"github.com/DownerCase/ecal-go/ecal"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelConfig struct {
	viewport viewport.Model
}

func NewConfigModel() *ModelConfig {
	viewport := viewport.New(85, 10)
	viewport.SetContent(ecal.GetConfig())
	viewport.Style = baseStyle

	return &ModelConfig{
		viewport: viewport,
	}
}

func (m *ModelConfig) Refresh() {}

func (m *ModelConfig) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	return cmd
}

func (m *ModelConfig) View() string {
	return m.viewport.View()
}
