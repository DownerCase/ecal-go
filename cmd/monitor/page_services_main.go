package main

import (
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model_services_main struct {
	table_services table.Model
}

func NewServicesMainModel() *model_services_main {
	cols := []table.Column{
		{Title: "ID", Width: 0}, // Hidden unique ID
		{Title: "T", Width: 1},  // Type (Client/Server)
		{Title: "Service", Width: 40},
		{Title: "Unit", Width: 32},
		{Title: "Tick", Width: 4},
	}

	return &model_services_main{
		table_services: NewTable(cols),
	}
}

func (m *model_services_main) Refresh() {
	m.updateTable(nil)
}

func (m *model_services_main) Update(msg tea.Msg) tea.Cmd {
	return m.updateTable(msg)
}

func (m *model_services_main) View() string {
	return baseStyle.Render(m.table_services.View()) + "\n" + m.table_services.HelpView()
}

func (m *model_services_main) GetSelectedId() (string, bool, error) {
	row := m.table_services.SelectedRow()
	if row == nil {
		return "", false, errEmptyTable
	}
	return row[0], row[1] == "S", nil
}

func (m *model_services_main) updateTable(msg tea.Msg) (cmd tea.Cmd) {
	rows := []table.Row{}
	mon := monitoring.GetMonitoring(monitoring.MonitorClient | monitoring.MonitorServer)
	for _, client := range mon.Clients {
		rows = append(rows, clientToRow(client))
	}
	for _, server := range mon.Servers {
		rows = append(rows, serverToRow(server))
	}
	m.table_services.SetRows(rows)
	m.table_services, cmd = m.table_services.Update(msg)
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
