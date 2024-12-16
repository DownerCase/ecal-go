package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ServicesPage int

const (
	subpage_services_main ServicesPage = iota
	subpage_services_detailed
)

type model_services struct {
	subpage ServicesPage
	pages   map[ServicesPage]PageModel
	NavKeys NavKeyMap
}

func NewServicesModel() *model_services {
	return (&model_services{
		subpage: subpage_services_main,
		pages: map[ServicesPage]PageModel{
			subpage_services_main:     NewServicesMainModel(),
			subpage_services_detailed: NewDetailedServiceModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *model_services) Refresh() {
	m.pages[m.subpage].Refresh()
}

func (m *model_services) Init() *model_services {
	m.NavKeys["esc"] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys["enter"] = func() tea.Cmd { m.navDown(); return nil }
	return m
}

func (m *model_services) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}
	return m.pages[m.subpage].Update(msg)
}

func (m *model_services) View() string {
	return m.pages[m.subpage].View()
}

func (m *model_services) navDown() {
	switch m.subpage {
	case subpage_services_main:
		main := m.pages[subpage_services_main].(*model_services_main)
		id, isServer, err := main.GetSelectedId()
		if err != nil {
			return // Can't transition
		}
		detailed := m.pages[subpage_services_detailed].(*model_service_detailed)
		detailed.IsServer = isServer
		detailed.Id = id
		detailed.Refresh()
		m.subpage = subpage_services_detailed

	}
}

func (m *model_services) navUp() {
	switch m.subpage {
	case subpage_services_detailed:
		m.subpage = subpage_services_main
	}
}
