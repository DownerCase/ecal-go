package main

import (
	"path/filepath"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model_processes struct {
	table_processes table.Model
}

func NewProcessesModel() model_processes {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "Name", Width: 33},
		{Title: "State", Width: 8},
		{Title: "Info", Width: 23},
		{Title: "Tick", Width: 4},
	}

	t := table.New(
		table.WithHeight(8),
		table.WithFocused(true),
		table.WithColumns(columns),
		table.WithStyles(tableStyle),
	)
	return model_processes{
		table_processes: t,
	}
}

func (m *model_processes) Init() tea.Cmd {
	return nil
}

func (m *model_processes) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TickMsg:
		m.Refresh()
		return doTick()
	case tea.KeyMsg:
		m.updateTable(msg)
	}
	return cmd
}

func (m *model_processes) View() string {
	return baseStyle.Render(m.table_processes.View()) + "\n" + m.table_processes.HelpView()
}

func (m *model_processes) Refresh() {
	m.updateTable(nil)
}

func (m *model_processes) updateTable(msg tea.Msg) {
	rows := []table.Row{}
	mon := monitoring.GetMonitoring(monitoring.MonitorProcess)
	for _, proc := range mon.Processes {
		rows = append(rows, procToRow(proc))
	}
	m.table_processes.SetRows(rows)
	m.table_processes, _ = m.table_processes.Update(msg)
}

func procToRow(proc monitoring.ProcessMon) table.Row {
	return []string{
		strconv.FormatInt(int64(proc.Pid), 10),
		filepath.Base(proc.Process_name),
		proc.State_severity.String(),
		proc.State_info,
		strconv.FormatInt(int64(proc.Registration_clock), 10),
	}
}
