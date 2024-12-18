package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
)

type modelHosts struct {
	table table.Model
}

func NewHostsModel() *modelHosts {
	columns := []table.Column{
		{Title: "Host", Width: 28},
		{Title: "Processes", Width: 9},
		{Title: "Subscribers", Width: 11},
		{Title: "Publishers", Width: 11},
		{Title: "Servers", Width: 7},
		{Title: "Clients", Width: 7},
	}

	return &modelHosts{
		table: NewTable(columns),
	}
}

func (m *modelHosts) Update(msg tea.Msg) tea.Cmd {
	m.updateTable(nil)
	return nil
}

func (m *modelHosts) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *modelHosts) Refresh() {
	m.updateTable(nil)
}

type hostInfo struct {
	Publishers  int
	Subscribers int
	Clients     int
	Servers     int
	Processes   int
}

func (m *modelHosts) updateTable(msg tea.Msg) {
	mon := monitoring.GetMonitoring(monitoring.MonitorAll)

	hosts := make(map[string]hostInfo)
	for _, pub := range mon.Publishers {
		host := hosts[pub.HostName]
		host.Publishers += 1
		hosts[pub.HostName] = host
	}
	for _, sub := range mon.Subscribers {
		host := hosts[sub.HostName]
		host.Subscribers += 1
		hosts[sub.HostName] = host
	}
	for _, client := range mon.Clients {
		host := hosts[client.HostName]
		host.Clients += 1
		hosts[client.HostName] = host
	}
	for _, server := range mon.Servers {
		host := hosts[server.HostName]
		host.Servers += 1
		hosts[server.HostName] = host
	}
	for _, proc := range mon.Processes {
		host := hosts[proc.HostName]
		host.Processes += 1
		hosts[proc.HostName] = host
	}

	m.table.SetRows(hostsToRows(hosts))
	m.table, _ = m.table.Update(msg)
}

func hostsToRows(hosts map[string]hostInfo) (rows []table.Row) {
	for hostName, hostInfo := range hosts {
		rows = append(rows, table.Row{
			hostName,
			strconv.FormatInt(int64(hostInfo.Processes), 10),
			strconv.FormatInt(int64(hostInfo.Subscribers), 10),
			strconv.FormatInt(int64(hostInfo.Publishers), 10),
			strconv.FormatInt(int64(hostInfo.Servers), 10),
			strconv.FormatInt(int64(hostInfo.Clients), 10),
		})
	}
	return
}
