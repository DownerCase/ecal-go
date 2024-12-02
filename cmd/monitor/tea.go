package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var highlight = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))

var p *tea.Program

type Page int

const (
	page_topics Page = iota
	page_services
	page_hosts
	page_processes
	page_logs
	page_system
	page_about
)

type model struct {
	page            Page
	model_topics    model_topics
	model_processes model_processes
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *model) updatePage(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.page {
	case page_topics:
		cmd = m.model_topics.Update(msg)
	case page_processes:
		cmd = m.model_processes.Update(msg)
	}
	return cmd
}

func (m *model) transitionTo(newPage Page) {
	switch newPage {
	case page_topics:
		m.model_topics.Refresh()
	case page_processes:
		m.model_processes.Refresh()
	}
	m.page = newPage
}

func (m *model) Init() tea.Cmd {
	return doTick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.transitionTo(page_topics)
		case "2":
			m.transitionTo(page_services)
		case "3":
			m.transitionTo(page_hosts)
		case "4":
			m.transitionTo(page_processes)
		case "5":
			m.transitionTo(page_logs)
		case "6":
			m.transitionTo(page_system)
		case "7":
			m.transitionTo(page_about)
		default:
			cmd = m.updatePage(msg)
		}
	default:
		cmd = m.updatePage(msg)
	}
	return m, cmd
}

func (m *model) View() string {
	s := strings.Builder{}
	switch m.page {
	case page_topics:
		s.WriteString(m.model_topics.View())
	case page_processes:
		s.WriteString(m.model_processes.View())
	default:
		s.WriteString(placeholderTab(m.page))
	}
	tabs := []string{
		"1: Topics",
		"2: Services",
		"3: Hosts",
		"4: Processes",
		"5: Logs",
		"6: System",
		"7: About",
	}
	s.WriteString("\n")
	page := m.page
	tabs[page] = highlight.Render(tabs[page])
	for _, tab := range tabs {
		s.WriteString(tab)
		s.WriteRune(' ')
	}
	return s.String()
}

func doCli() {

	m := model{page_topics, NewTopicsModel(), NewProcessesModel()}
	p = tea.NewProgram(&m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
