package main

import (
	"fmt"
	"strings"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelConfig struct {
	viewport viewport.Model
}

func getContent() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%-24s : %s (%s)\n", "eCAL Version", ecal.GetVersionString(), ecal.GetVersionDateString()))

	if !ecal.IsInitialized() {
		s.WriteString("eCAL not initialized!\n")
		return s.String()
	}

	s.WriteString(fmt.Sprintf("%-24s : %s \n", "Config", ecal.GetLoadedConfigFilePath()))

	s.WriteString("\n--------------------------Publisher--------------------------\n")
	s.WriteString(fmt.Sprintf("%-24s : ", "Layers"))

	{
		publisherLayers := []string{}
		if ecal.PublisherShmEnabled() {
			publisherLayers = append(publisherLayers, "SHM")
		}

		if ecal.PublisherUDPEnabled() {
			publisherLayers = append(publisherLayers, "UDP")
		}

		if ecal.PublisherTCPEnabled() {
			publisherLayers = append(publisherLayers, "TCP")
		}

		s.WriteString(fmt.Sprintln(publisherLayers))
	}

	s.WriteString("\n--------------------------Subscriber-------------------------\n")
	s.WriteString(fmt.Sprintf("%-24s : ", "Layers"))

	{
		publisherLayers := []string{}
		if ecal.SubscriberShmEnabled() {
			publisherLayers = append(publisherLayers, "SHM")
		}

		if ecal.SubscriberUDPEnabled() {
			publisherLayers = append(publisherLayers, "UDP")
		}

		if ecal.SubscriberTCPEnabled() {
			publisherLayers = append(publisherLayers, "TCP")
		}

		s.WriteString(fmt.Sprintln(publisherLayers))
	}

	return s.String()
}

func NewConfigModel() *ModelConfig {
	viewport := viewport.New(85, 10)

	viewport.SetContent(getContent())
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
