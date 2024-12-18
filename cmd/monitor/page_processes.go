package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ProcessesPage int

const (
	subpageProcMain ProcessesPage = iota
	subpageProcDetailed
)

type modelProcesses struct {
	subpage ProcessesPage
	pages   map[ProcessesPage]PageModel
	NavKeys NavKeyMap
}

func NewProcessesModel() *modelProcesses {
	return (&modelProcesses{
		subpage: subpageProcMain,
		pages: map[ProcessesPage]PageModel{
			subpageProcMain:     NewProcessesMainModel(),
			subpageProcDetailed: NewDetailedProcessModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *modelProcesses) Init() *modelProcesses {
	m.NavKeys["esc"] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys["enter"] = func() tea.Cmd { m.navDown(); return nil }
	return m
}

func (m *modelProcesses) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}
	return m.pages[m.subpage].Update(msg)
}

func (m *modelProcesses) View() string {
	return m.pages[m.subpage].View()
}

func (m *modelProcesses) Refresh() {
	m.pages[m.subpage].Refresh()
}

func (m *modelProcesses) navDown() {
	if m.subpage == subpageProcMain {
		main := m.pages[subpageProcMain].(*modelProcessesMain)
		pid, err := main.getSelectedPid()
		if err != nil {
			return // Can't transition
		}
		detailed := m.pages[subpageProcDetailed].(*modelProcessDetailed)
		detailed.Pid = pid
		m.subpage = subpageProcDetailed
		detailed.Refresh()
	}
}

func (m *modelProcesses) navUp() {
	if m.subpage == subpageProcDetailed {
		m.subpage = subpageProcMain
	}
}
