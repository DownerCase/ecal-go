package main

import (
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type modelServicesMain struct {
	table table.Model
}

func NewServicesMainModel() *modelServicesMain {
	cols := []table.Column{
		{Title: "ID", Width: 0}, // Hidden unique ID
		{Title: "T", Width: 1},  // Type (Client/Server)
		{Title: "Service", Width: 40},
		{Title: "Unit", Width: 32},
		{Title: "Tick", Width: 4},
	}

	return &modelServicesMain{
		table: NewTable(cols),
	}
}

func (m *modelServicesMain) Refresh() {
	m.updateTable(nil)
}

func (m *modelServicesMain) Update(msg tea.Msg) tea.Cmd {
	return m.updateTable(msg)
}

func (m *modelServicesMain) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *modelServicesMain) GetSelectedID() (string, bool, error) {
	row := m.table.SelectedRow()
	if row == nil {
		return "", false, errEmptyTable
	}
	return row[0], row[1] == "S", nil
}

func (m *modelServicesMain) updateTable(msg tea.Msg) (cmd tea.Cmd) {
	rows := []table.Row{}
	mon := monitoring.GetMonitoring(monitoring.MonitorClient | monitoring.MonitorServer)
	for _, client := range mon.Clients {
		rows = append(rows, clientToRow(client))
	}
	for _, server := range mon.Servers {
		rows = append(rows, serverToRow(server))
	}
	m.table.SetRows(rows)
	m.table, cmd = m.table.Update(msg)
	return
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
		[]string{client.ID, "C"},
		serviceToRow(client.ServiceBase)...,
	)
}

func serverToRow(server monitoring.ServerMon) table.Row {
	return append(
		[]string{server.ID, "S"},
		serviceToRow(server.ServiceBase)...,
	)
}
