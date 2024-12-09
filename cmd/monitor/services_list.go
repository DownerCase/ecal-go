package main

import (
	"errors"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ServicesPage int

const (
	subpage_services_main ServicesPage = iota
	subpage_services_detailed
)

type model_services struct {
	table_services table.Model
	subpage        ServicesPage
	model_detailed model_service_detailed
	NavKeys        NavKeyMap
}

func NewServicesModel() *model_services {
	cols := []table.Column{
		{Title: "ID", Width: 0}, // Hidden unique ID
		{Title: "T", Width: 1},  // Type (Client/Server)
		{Title: "Service", Width: 40},
		{Title: "Unit", Width: 32},
		{Title: "Tick", Width: 4},
	}

	return (&model_services{
		table_services: NewTable(cols),
		subpage:        subpage_services_main,
		model_detailed: NewDetailedServiceModel(),
		NavKeys:        make(NavKeyMap),
	}).Init()
}

func (m *model_services) Refresh() {
	switch m.subpage {
	case subpage_services_detailed:
		m.model_detailed.Refresh()
	default:
		m.updateTable(nil)
	}
}

func (m *model_services) Init() *model_services {
	m.NavKeys[tea.KeyEscape] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys[tea.KeyEnter] = func() tea.Cmd { m.navDown(); return nil }
	return m
}

func (m *model_services) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}

	var cmd tea.Cmd
	switch m.subpage {
	case subpage_services_main:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			m.updateTable(msg)
		}
	case subpage_services_detailed:
		cmd = m.model_detailed.Update(msg)
	}
	return cmd
}

func (m *model_services) View() string {
	switch m.subpage {
	case subpage_services_main:
		return baseStyle.Render(m.table_services.View()) + "\n" + m.table_services.HelpView()
	case subpage_services_detailed:
		return m.model_detailed.View()
	}
	return "Invalid page"
}

func (m *model_services) navDown() {
	switch m.subpage {
	case subpage_services_main:
		id, isServer, err := m.getSelectedId()
		if err != nil {
			return // Can't transition
		}
		m.model_detailed.Id = id
		m.model_detailed.IsServer = isServer
		m.subpage = subpage_services_detailed
		m.model_detailed.Refresh()
	}
}

func (m *model_services) navUp() {
	switch m.subpage {
	case subpage_services_detailed:
		m.subpage = subpage_services_main
	}
}

func (m *model_services) getSelectedId() (string, bool, error) {
	row := m.table_services.SelectedRow()
	if row == nil {
		return "", false, errors.New("No processes")
	}
	return row[0], row[1] == "S", nil
}

func (m *model_services) updateTable(msg tea.Msg) {
	rows := []table.Row{}
	mon := monitoring.GetMonitoring(monitoring.MonitorClient | monitoring.MonitorServer)
	for _, client := range mon.Clients {
		rows = append(rows, clientToRow(client))
	}
	for _, server := range mon.Servers {
		rows = append(rows, serverToRow(server))
	}
	m.table_services.SetRows(rows)
	m.table_services, _ = m.table_services.Update(msg)
}

func serviceToRow(service monitoring.ServiceBase) table.Row {
	return []string{
		service.Name,
		service.Unit,
		strconv.FormatInt(int64(service.RegistrationClock), 10),
	}
}

func clientToRow(client monitoring.ClientMon) table.Row {
	return append(
		[]string{client.Id, "C"},
		serviceToRow(client.ServiceBase)...,
	)
}

func serverToRow(server monitoring.ServerMon) table.Row {
	return append(
		[]string{server.Id, "S"},
		serviceToRow(server.ServiceBase)...,
	)
}
