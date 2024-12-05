package main

import (
	"fmt"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model_service_detailed struct {
	table_detailed table.Model
	Id             string
	IsServer       bool
}

func NewDetailedServiceModel() model_service_detailed {
	cols := []table.Column{
		{Title: "", Width: 10},
		{Title: "", Width: 67},
	}
	t := table.New(
		table.WithHeight(8),
		table.WithFocused(true),
		table.WithColumns(cols),
		table.WithStyles(tableStyle),
	)

	return model_service_detailed{
		table_detailed: t,
		Id:             "",
	}
}

func (m *model_service_detailed) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.table_detailed, cmd = m.table_detailed.Update(msg)
	}
	return cmd
}

func (m *model_service_detailed) View() string {
	return baseStyle.Render(m.table_detailed.View()) + "\n" + m.table_detailed.HelpView()
}

func (m *model_service_detailed) Refresh() {
	m.updateDetailedTable(nil)
}

func (m *model_service_detailed) updateDetailedTable(msg tea.Msg) {
	if m.IsServer {
		mon := monitoring.GetMonitoring(monitoring.MonitorServer)
		for _, s := range mon.Servers {
			if s.Id == m.Id {
				m.updateTableServer(&s)
				break
			}
		}
	} else {
		mon := monitoring.GetMonitoring(monitoring.MonitorClient)
		for _, c := range mon.Clients {
			if c.Id == m.Id {
				m.updateTableClient(&c)
				break
			}
		}
	}
	m.table_detailed, _ = m.table_detailed.Update(msg)
}

func (m *model_service_detailed) updateTableClient(c *monitoring.ClientMon) {
	m.table_detailed.Columns()[0].Title = "Client"
	m.table_detailed.Columns()[1].Title = c.Process
	baseRows := m.getBaseRows(c.ServiceBase)
	m.table_detailed.SetRows(append(baseRows, getMethodRows(c.ServiceBase)...))
}

func (m *model_service_detailed) updateTableServer(s *monitoring.ServerMon) {
	m.table_detailed.Columns()[0].Title = "Server"
	m.table_detailed.Columns()[1].Title = s.Process
	baseRows := m.getBaseRows(s.ServiceBase)
	baseRows = append(baseRows, table.Row{"TCP Port", fmt.Sprintf("V0: %v, V1: %v", s.PortV0, s.PortV1)})
	m.table_detailed.SetRows(append(baseRows, getMethodRows(s.ServiceBase)...))
}

func (m *model_service_detailed) getBaseRows(b monitoring.ServiceBase) []table.Row {
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
