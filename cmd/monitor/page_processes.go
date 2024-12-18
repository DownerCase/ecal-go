package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ProcessesPage int

const (
	subpageProcMain ProcessesPage = iota
	subpageProcDetailed
)

type ModelProcesses struct {
	subpage ProcessesPage
	pages   map[ProcessesPage]PageModel
	NavKeys NavKeyMap
}

func NewProcessesModel() *ModelProcesses {
	return (&ModelProcesses{
		subpage: subpageProcMain,
		pages: map[ProcessesPage]PageModel{
			subpageProcMain:     NewProcessesMainModel(),
			subpageProcDetailed: NewDetailedProcessModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *ModelProcesses) Init() *ModelProcesses {
	m.NavKeys["esc"] = func() tea.Cmd { return m.navUp() }
	m.NavKeys["enter"] = func() tea.Cmd { return m.navDown() }
	return m
}

func (m *ModelProcesses) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}
	return m.pages[m.subpage].Update(msg)
}

func (m *ModelProcesses) View() string {
	return m.pages[m.subpage].View()
}

func (m *ModelProcesses) Refresh() {
	m.pages[m.subpage].Refresh()
}

func (m *ModelProcesses) navDown() tea.Cmd {
	if m.subpage == subpageProcMain {
		main := m.pages[subpageProcMain].(*ModelProcessesMain)
		pid, err := main.getSelectedPid()
		if err != nil {
			return nil // Can't transition
		}
		detailed := m.pages[subpageProcDetailed].(*ModelProcessDetailed)
		detailed.Pid = pid
		m.subpage = subpageProcDetailed
		detailed.Refresh()
	}
	return nil
}

func (m *ModelProcesses) navUp() tea.Cmd {
	if m.subpage == subpageProcDetailed {
		m.subpage = subpageProcMain
	}
	return nil
}
