package main

import (
	"errors"
	"path/filepath"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model_processes_main struct {
	table_processes table.Model
}

func NewProcessesMainModel() *model_processes_main {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "Name", Width: 33},
		{Title: "State", Width: 8},
		{Title: "Info", Width: 23},
		{Title: "Tick", Width: 4},
	}

	return &model_processes_main{
		table_processes: NewTable(columns),
	}
}

func (m *model_processes_main) Update(msg tea.Msg) tea.Cmd {
	return m.updateTable(msg)
}

func (m *model_processes_main) View() string {
	return baseStyle.Render(m.table_processes.View()) + "\n" + m.table_processes.HelpView()
}

func (m *model_processes_main) Refresh() {
	m.updateTable(nil)
}

func (m *model_processes_main) getSelectedPid() (int32, error) {
	row := m.table_processes.SelectedRow()
	if row == nil {
		return 0, errors.New("No processes")
	}
	pid, err := strconv.ParseInt(row[0], 10, 64)
	return int32(pid), err
}

func (m *model_processes_main) updateTable(msg tea.Msg) (cmd tea.Cmd) {
	rows := []table.Row{}
	mon := monitoring.GetMonitoring(monitoring.MonitorProcess)
	for _, proc := range mon.Processes {
		rows = append(rows, procToRow(proc))
	}
	m.table_processes.SetRows(rows)
	m.table_processes, cmd = m.table_processes.Update(msg)
	return
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
