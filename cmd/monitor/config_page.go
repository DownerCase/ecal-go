package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/DownerCase/ecal-go/ecal"
)

type modelConfig struct {
	viewport viewport.Model
}

func NewConfigModel() *modelConfig {
	viewport := viewport.New(85, 10)
	viewport.SetContent(ecal.GetConfig())
	viewport.Style = baseStyle
	return &modelConfig{
		viewport: viewport,
	}
}

func (m *modelConfig) Refresh() {}

func (m *modelConfig) Update(msg tea.Msg) (cmd tea.Cmd) {
	m.viewport, cmd = m.viewport.Update(msg)
	return cmd
}

func (m *modelConfig) View() string {
	return m.viewport.View()
}
