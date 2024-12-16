package main

import (
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type PageModel interface {
	Refresh()
	Update(tea.Msg) tea.Cmd
	View() string
}

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
	page  Page
	pages map[Page]PageModel
}

func newModel() *model {
	pagesMap := make(map[Page]PageModel)
	pagesMap[page_topics] = NewTopicsModel()
	pagesMap[page_services] = NewServicesModel()
	pagesMap[page_hosts] = NewHostsModel()
	pagesMap[page_processes] = NewProcessesModel()
	pagesMap[page_logs] = NewLogsModel()
	pagesMap[page_system] = NewConfigModel()
	pagesMap[page_about] = &PlaceholderModel{"About Placeholder"}
	return &model{
		page:  page_topics,
		pages: pagesMap,
	}
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *model) updatePage(msg tea.Msg) tea.Cmd {
	return m.pages[m.page].Update(msg)
}

func (m *model) transitionTo(newPage Page) {
	m.page = newPage
	m.refresh()
}

func (m *model) refresh() {
	m.pages[m.page].Refresh()
}

func (m *model) Init() tea.Cmd {
	return doTick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TickMsg:
		m.refresh()
		cmd = doTick()
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
	s.WriteString(m.pages[m.page].View())
	tabs := []string{
		"1: Topics",
		"2: Services",
		"3: Hosts",
		"4: Processes",
		"5: Logs",
		"6: Config",
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
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
