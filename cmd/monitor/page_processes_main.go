package main

import (
	"path/filepath"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelProcessesMain struct {
	table table.Model
}

func NewProcessesMainModel() *ModelProcessesMain {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "Name", Width: 33},
		{Title: "State", Width: 8},
		{Title: "Info", Width: 23},
		{Title: "Tick", Width: 4},
	}

	return &ModelProcessesMain{
		table: NewTable(columns),
	}
}

func (m *ModelProcessesMain) Update(msg tea.Msg) tea.Cmd {
	return m.updateTable(msg)
}

func (m *ModelProcessesMain) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *ModelProcessesMain) Refresh() {
	m.updateTable(nil)
}

func (m *ModelProcessesMain) getSelectedPid() (int32, error) {
	row := m.table.SelectedRow()
	if row == nil {
		return 0, errEmptyTable
	}

	pid, err := strconv.ParseInt(row[0], 10, 32)

	return int32(pid), err
}

func (m *ModelProcessesMain) updateTable(msg tea.Msg) tea.Cmd {
	rows := []table.Row{}

	mon := monitoring.GetMonitoring(monitoring.MonitorProcess)
	for _, proc := range mon.Processes {
		rows = append(rows, procToRow(proc))
	}

	var cmd tea.Cmd

	m.table.SetRows(rows)
	m.table, cmd = m.table.Update(msg)

	return cmd
}

func procToRow(proc monitoring.ProcessMon) table.Row {
	return []string{
		strconv.FormatInt(int64(proc.Pid), 10),
		filepath.Base(proc.ProcessName),
		proc.StateSeverity.String(),
		proc.StateInfo,
		strconv.FormatInt(int64(proc.RegistrationClock), 10),
	}
}
