package main

import (
	"fmt"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model_process_detailed struct {
	table_detailed table.Model
	Pid            int32
}

func NewDetailedProcessModel() model_process_detailed {
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

	return model_process_detailed{
		table_detailed: t,
		Pid:            0,
	}
}

func (m *model_process_detailed) Init() tea.Cmd {
	return nil
}

func (m *model_process_detailed) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TickMsg:
		m.updateDetailedTable(nil)
		return doTick()
	case tea.KeyMsg:
		m.table_detailed, cmd = m.table_detailed.Update(msg)
	}
	return cmd
}

func (m *model_process_detailed) View() string {
	return baseStyle.Render(m.table_detailed.View()) + "\n" + m.table_detailed.HelpView()
}

func (m *model_process_detailed) Refresh() {
	m.updateDetailedTable(nil)
}

func (m *model_process_detailed) updateDetailedTable(msg tea.Msg) {
	mon := monitoring.GetMonitoring(monitoring.MonitorProcess)
	var p monitoring.ProcessMon

	for _, proc := range mon.Processes {
		if proc.Pid == m.Pid {
			p = proc
			break
		}
	}
	m.table_detailed.Columns()[0].Title = strconv.FormatInt(int64(p.Pid), 10)
	m.table_detailed.Columns()[1].Title = p.Process_name
	health := fmt.Sprintf("%s %v", p.State_severity.String(), p.State_severity_level)
	rows := []table.Row{
		{"Unit", p.Unit_name},
		{health, p.State_info},
		{"Parameters", p.Process_parameters},
		{"Host", fmt.Sprintf("%s (Group: %s)", p.Host_name, p.Host_group)},
		{"Components", p.Components_initialized},
		{"Runtime", p.Runtime_version},
		{"Tick", strconv.FormatInt(int64(p.Registration_clock), 10)},
	}
	m.table_detailed.SetRows(rows)
	m.table_detailed, _ = m.table_detailed.Update(msg)
}
