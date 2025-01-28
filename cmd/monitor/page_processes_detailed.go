package main

import (
	"fmt"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelProcessDetailed struct {
	table table.Model
	Pid   int32
}

func NewDetailedProcessModel() *ModelProcessDetailed {
	cols := []table.Column{
		{Title: "", Width: 10},
		{Title: "", Width: 67},
	}

	return &ModelProcessDetailed{
		table: NewTable(cols),
		Pid:   0,
	}
}

func (m *ModelProcessDetailed) Init() tea.Cmd {
	return nil
}

func (m *ModelProcessDetailed) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
		m.table, cmd = m.table.Update(msg)
	}

	return cmd
}

func (m *ModelProcessDetailed) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *ModelProcessDetailed) Refresh() {
	m.updateDetailedTable(nil)
}

func (m *ModelProcessDetailed) updateDetailedTable(msg tea.Msg) {
	var p monitoring.ProcessMon

	mon := monitoring.GetMonitoring(monitoring.MonitorProcess)
	for _, proc := range mon.Processes {
		if proc.Pid == m.Pid {
			p = proc
			break
		}
	}

	m.table.Columns()[0].Title = strconv.FormatInt(int64(p.Pid), 10)
	m.table.Columns()[1].Title = p.ProcessName
	health := fmt.Sprintf("%s %v", p.StateSeverity.String(), p.StateSeverityLevel)
	rows := []table.Row{
		{"Unit", p.UnitName},
		{health, p.StateInfo},
		{"Parameters", p.ProcessParameters},
		{"Host", fmt.Sprintf("%s (SHM Domain: %s)", p.HostName, p.ShmDomain)},
		{"Components", p.ComponentsInitialized},
		{"Runtime", p.RuntimeVersion},
		{"Tick", strconv.FormatInt(int64(p.RegistrationClock), 10)},
	}
	m.table.SetRows(rows)
	m.table, _ = m.table.Update(msg)
}
