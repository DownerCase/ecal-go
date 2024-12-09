package main

import (
	"github.com/DownerCase/ecal-go/ecal"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model_config struct {
	viewport viewport.Model
}

func NewConfigModel() *model_config {
	viewport := viewport.New(85, 10)
	viewport.SetContent(ecal.GetConfig())
	viewport.Style = baseStyle
	return &model_config{
		viewport: viewport,
	}
}

func (m *model_config) Refresh() {}

func (m *model_config) Update(msg tea.Msg) (cmd tea.Cmd) {
	m.viewport, cmd = m.viewport.Update(msg)
	return cmd
}

func (m *model_config) View() string {
	return m.viewport.View()
}
