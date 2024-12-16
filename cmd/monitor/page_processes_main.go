package main

import (
	"path/filepath"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type modelProcessesMain struct {
	table table.Model
}

func NewProcessesMainModel() *modelProcessesMain {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "Name", Width: 33},
		{Title: "State", Width: 8},
		{Title: "Info", Width: 23},
		{Title: "Tick", Width: 4},
	}

	return &modelProcessesMain{
		table: NewTable(columns),
	}
}

func (m *modelProcessesMain) Update(msg tea.Msg) tea.Cmd {
	return m.updateTable(msg)
}

func (m *modelProcessesMain) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *modelProcessesMain) Refresh() {
	m.updateTable(nil)
}

func (m *modelProcessesMain) getSelectedPid() (int32, error) {
	row := m.table.SelectedRow()
	if row == nil {
		return 0, errEmptyTable
	}
	pid, err := strconv.ParseInt(row[0], 10, 64)
	return int32(pid), err
}

func (m *modelProcessesMain) updateTable(msg tea.Msg) (cmd tea.Cmd) {
	rows := []table.Row{}
	mon := monitoring.GetMonitoring(monitoring.MonitorProcess)
	for _, proc := range mon.Processes {
		rows = append(rows, procToRow(proc))
	}
	m.table.SetRows(rows)
	m.table, cmd = m.table.Update(msg)
	return
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
