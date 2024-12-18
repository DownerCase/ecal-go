package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ServicesPage int

const (
	subpageServicesMain ServicesPage = iota
	subpageServicesDetailed
)

type ModelServices struct {
	subpage ServicesPage
	pages   map[ServicesPage]PageModel
	NavKeys NavKeyMap
}

func NewServicesModel() *ModelServices {
	return (&ModelServices{
		subpage: subpageServicesMain,
		pages: map[ServicesPage]PageModel{
			subpageServicesMain:     NewServicesMainModel(),
			subpageServicesDetailed: NewDetailedServiceModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *ModelServices) Refresh() {
	m.pages[m.subpage].Refresh()
}

func (m *ModelServices) Init() *ModelServices {
	m.NavKeys["esc"] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys["enter"] = func() tea.Cmd { m.navDown(); return nil }
	return m
}

func (m *ModelServices) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}
	return m.pages[m.subpage].Update(msg)
}

func (m *ModelServices) View() string {
	return m.pages[m.subpage].View()
}

func (m *ModelServices) navDown() {
	if m.subpage == subpageServicesMain {
		main := m.pages[subpageServicesMain].(*ModelServicesMain)
		id, isServer, err := main.GetSelectedID()
		if err != nil {
			return // Can't transition
		}
		detailed := m.pages[subpageServicesDetailed].(*ModelServiceDetailed)
		detailed.IsServer = isServer
		detailed.ID = id
		detailed.Refresh()
		m.subpage = subpageServicesDetailed
	}
}

func (m *ModelServices) navUp() {
	if m.subpage == subpageServicesDetailed {
		m.subpage = subpageServicesMain
	}
}
