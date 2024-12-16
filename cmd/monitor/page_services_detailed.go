package main

import (
	"fmt"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type modelServiceDetailed struct {
	table    table.Model
	ID       string
	IsServer bool
}

func NewDetailedServiceModel() *modelServiceDetailed {
	cols := []table.Column{
		{Title: "", Width: 10},
		{Title: "", Width: 67},
	}

	return &modelServiceDetailed{
		table: NewTable(cols),
		ID:    "",
	}
}

func (m *modelServiceDetailed) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.table, cmd = m.table.Update(msg)
	}
	return cmd
}

func (m *modelServiceDetailed) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *modelServiceDetailed) Refresh() {
	m.updateDetailedTable(nil)
}

func (m *modelServiceDetailed) updateDetailedTable(msg tea.Msg) {
	if m.IsServer {
		mon := monitoring.GetMonitoring(monitoring.MonitorServer)
		for _, s := range mon.Servers {
			if s.ID == m.ID {
				m.updateTableServer(&s)
				break
			}
		}
	} else {
		mon := monitoring.GetMonitoring(monitoring.MonitorClient)
		for _, c := range mon.Clients {
			if c.ID == m.ID {
				m.updateTableClient(&c)
				break
			}
		}
	}
	m.table, _ = m.table.Update(msg)
}

func (m *modelServiceDetailed) updateTableClient(c *monitoring.ClientMon) {
	m.table.Columns()[0].Title = "Client"
	m.table.Columns()[1].Title = c.Process
	baseRows := m.getBaseRows(c.ServiceBase)
	m.table.SetRows(append(baseRows, getMethodRows(c.ServiceBase)...))
}

func (m *modelServiceDetailed) updateTableServer(s *monitoring.ServerMon) {
	m.table.Columns()[0].Title = "Server"
	m.table.Columns()[1].Title = s.Process
	baseRows := m.getBaseRows(s.ServiceBase)
	baseRows = append(baseRows, table.Row{"TCP Port", fmt.Sprintf("V0: %v, V1: %v", s.PortV0, s.PortV1)})
	m.table.SetRows(append(baseRows, getMethodRows(s.ServiceBase)...))
}

func (m *modelServiceDetailed) getBaseRows(b monitoring.ServiceBase) []table.Row {
	return []table.Row{
		{"Unit", fmt.Sprintf("%s (Protocol V%v)", b.Unit, b.ProtocolVersion)},
		{"Pid", fmt.Sprintf("%v (%s)", b.Pid, b.HostName)},
		{"Tick", strconv.FormatInt(int64(b.RegistrationClock), 10)},
	}
}

func getMethodRows(b monitoring.ServiceBase) []table.Row {
	rows := []table.Row{
		{"Methods", strconv.FormatInt(int64(len(b.Methods)), 10)},
	}
	for _, method := range b.Methods {
		rows = append(rows, table.Row{
			method.Name,
			fmt.Sprintf("%s -> %s (Called x%v)", method.RequestType.Type, method.ResponseType.Type, method.CallCount),
		})
	}
	return rows
}
