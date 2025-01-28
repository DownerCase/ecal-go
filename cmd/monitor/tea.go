package main

import (
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type PageModel interface {
	Refresh()
	Update(msg tea.Msg) tea.Cmd
	View() string
}

type Page int

const (
	pageTopics Page = iota
	pageServices
	pageHosts
	pageProcesses
	pageLogs
	pageSystem
)

type model struct {
	page  Page
	pages map[Page]PageModel
}

func newModel() *model {
	pagesMap := make(map[Page]PageModel)
	pagesMap[pageTopics] = NewTopicsModel()
	pagesMap[pageServices] = NewServicesModel()
	pagesMap[pageHosts] = NewHostsModel()
	pagesMap[pageProcesses] = NewProcessesModel()
	pagesMap[pageLogs] = NewLogsModel()
	pagesMap[pageSystem] = NewConfigModel()

	return &model{
		page:  pageTopics,
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
			m.transitionTo(pageTopics)
		case "2":
			m.transitionTo(pageServices)
		case "3":
			m.transitionTo(pageHosts)
		case "4":
			m.transitionTo(pageProcesses)
		case "5":
			m.transitionTo(pageLogs)
		case "6":
			m.transitionTo(pageSystem)
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
	s.WriteString("\n")

	tabs := []string{
		"1: Topics",
		"2: Services",
		"3: Hosts",
		"4: Processes",
		"5: Logs",
		"6: Config",
	}
	tabs[m.page] = highlight.Render(tabs[m.page])

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
