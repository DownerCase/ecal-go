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
	m.NavKeys["esc"] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys["enter"] = func() tea.Cmd { m.navDown(); return nil }
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

func (m *ModelProcesses) navDown() {
	if m.subpage == subpageProcMain {
		main := m.pages[subpageProcMain].(*ModelProcessesMain)
		pid, err := main.getSelectedPid()
		if err != nil {
			return // Can't transition
		}
		detailed := m.pages[subpageProcDetailed].(*ModelProcessDetailed)
		detailed.Pid = pid
		m.subpage = subpageProcDetailed
		detailed.Refresh()
	}
}

func (m *ModelProcesses) navUp() {
	if m.subpage == subpageProcDetailed {
		m.subpage = subpageProcMain
	}
}
