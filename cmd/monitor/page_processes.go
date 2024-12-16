package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ProcessesPage int

const (
	subpage_proc_main ProcessesPage = iota
	subpage_proc_detailed
)

type model_processes struct {
	subpage ProcessesPage
	pages   map[ProcessesPage]PageModel
	NavKeys NavKeyMap
}

func NewProcessesModel() *model_processes {
	return (&model_processes{
		subpage: subpage_proc_main,
		pages: map[ProcessesPage]PageModel{
			subpage_proc_main:     NewProcessesMainModel(),
			subpage_proc_detailed: NewDetailedProcessModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *model_processes) Init() *model_processes {
	m.NavKeys["esc"] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys["enter"] = func() tea.Cmd { m.navDown(); return nil }
	return m
}

func (m *model_processes) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}
	return m.pages[m.subpage].Update(msg)
}

func (m *model_processes) View() string {
	return m.pages[m.subpage].View()
}

func (m *model_processes) Refresh() {
	m.pages[m.subpage].Refresh()
}

func (m *model_processes) navDown() {
	switch m.subpage {
	case subpage_proc_main:
		main := m.pages[subpage_proc_main].(*model_processes_main)
		pid, err := main.getSelectedPid()
		if err != nil {
			return // Can't transition
		}
		detailed := m.pages[subpage_proc_detailed].(*model_process_detailed)
		detailed.Pid = pid
		m.subpage = subpage_proc_detailed
		detailed.Refresh()
	}
}

func (m *model_processes) navUp() {
	switch m.subpage {
	case subpage_proc_detailed:
		m.subpage = subpage_proc_main
	}
}
