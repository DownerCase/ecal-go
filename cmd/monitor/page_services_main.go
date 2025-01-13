package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
)

type ModelServicesMain struct {
	table table.Model
}

func NewServicesMainModel() *ModelServicesMain {
	cols := []table.Column{
		{Title: "ID", Width: 0}, // Hidden unique ID
		{Title: "T", Width: 1},  // Type (Client/Server)
		{Title: "Service", Width: 40},
		{Title: "Unit", Width: 32},
		{Title: "Tick", Width: 4},
	}

	return &ModelServicesMain{
		table: NewTable(cols),
	}
}

func (m *ModelServicesMain) Refresh() {
	m.updateTable(nil)
}

func (m *ModelServicesMain) Update(msg tea.Msg) tea.Cmd {
	return m.updateTable(msg)
}

func (m *ModelServicesMain) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *ModelServicesMain) GetSelectedID() (uint64, bool, error) {
	row := m.table.SelectedRow()
	if row == nil {
		return 0, false, errEmptyTable
	}

	id, err := strconv.ParseUint(row[0], 10, 64)
	if err != nil {
		err = fmt.Errorf("services - GetSelectedID() %w", err)
	}

	return id, row[1] == "S", err
}

func (m *ModelServicesMain) updateTable(msg tea.Msg) tea.Cmd {
	rows := []table.Row{}

	mon := monitoring.GetMonitoring(monitoring.MonitorClient | monitoring.MonitorServer)
	for _, client := range mon.Clients {
		rows = append(rows, clientToRow(client))
	}

	for _, server := range mon.Servers {
		rows = append(rows, serverToRow(server))
	}

	var cmd tea.Cmd

	m.table.SetRows(rows)
	m.table, cmd = m.table.Update(msg)

	return cmd
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
		[]string{strconv.FormatUint(client.ID, 10), "C"},
		serviceToRow(client.ServiceBase)...,
	)
}

func serverToRow(server monitoring.ServerMon) table.Row {
	return append(
		[]string{strconv.FormatUint(server.ID, 10), "S"},
		serviceToRow(server.ServiceBase)...,
	)
}
