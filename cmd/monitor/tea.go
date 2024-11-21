package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Bold(true)

var p *tea.Program

type Page int

const (
	page_topics Page = iota
	page_topic_detailed
)

type model struct {
	page           Page
	model_topics   model_topics
	model_detailed model_detailed
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *model) refresh() {
	switch m.page {
	case page_topics:
		m.model_topics.Refresh()
	case page_topic_detailed:
		m.model_detailed.Refresh()
	}
}

func (m *model) updatePage(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.page {
	case page_topics:
		cmd = m.model_topics.Update(msg)
	case page_topic_detailed:
		cmd = m.model_detailed.Update(msg)
	}
	return cmd
}

// Navigate down - i.e.: Handle Enter keypress
func (m *model) navDown() {
	switch m.page {
	case page_topics:
		m.transitionTo(page_topic_detailed)
	}
}

// Navigate down - i.e.: Handle Esc keypress
func (m *model) navUp() {
	switch m.page {
	case page_topic_detailed:
		m.transitionTo(page_topics)
	}
}

func (m *model) transitionTo(newPage Page) {
	// Pre update
	switch newPage {
	case page_topics:
		m.model_topics.Refresh()
	case page_topic_detailed:
		topic, is_subscriber ,err := m.model_topics.GetSelectedId()
		if err != nil {
			return // Don't' transition
		}
		m.model_detailed.ShowTopic(topic, is_subscriber)
	}
	m.page = newPage
}

func (m *model) Init() tea.Cmd {
	return doTick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TickMsg:
		m.refresh()
		return m, doTick()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.navUp()
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.navDown()
		default:
			m.updatePage(msg)
		}
	default:
		m.updatePage(msg)
	}
	return m, cmd
}

func (m *model) View() string {
	switch m.page {
	case page_topics:
		return m.model_topics.View()
	case page_topic_detailed:
		return m.model_detailed.View()
	default:
		return "Invalid page"
	}
}

func doCli() {

	m := model{page_topics, NewTopicsModel(), NewDetailedModel()}
	p = tea.NewProgram(&m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
