package main

import (
	"errors"
	"path/filepath"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ProcessesPage int

const (
	subpage_proc_main ProcessesPage = iota
	subpage_proc_detailed
)

type model_processes struct {
	table_processes table.Model
	subpage         ProcessesPage
	model_detailed  model_process_detailed
}

func NewProcessesModel() *model_processes {
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
	return &model_processes{
		table_processes: t,
		subpage:         subpage_proc_main,
		model_detailed:  NewDetailedProcessModel(),
	}
}

func (m *model_processes) Init() tea.Cmd {
	return nil
}

func (m *model_processes) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Global navigation keys
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			m.navUp()
			return cmd
		case tea.KeyEnter:
			m.navDown()
			return cmd
		}
	}

	switch m.subpage {
	case subpage_proc_main:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			m.updateTable(msg)
		}
	case subpage_proc_detailed:
		cmd = m.model_detailed.Update(msg)
	}
	return cmd
}

func (m *model_processes) View() string {
	switch m.subpage {
	case subpage_proc_main:
		return baseStyle.Render(m.table_processes.View()) + "\n" + m.table_processes.HelpView()
	case subpage_proc_detailed:
		return m.model_detailed.View()
	}
	return "Invalid page"
}

func (m *model_processes) Refresh() {
	switch m.subpage {
	case subpage_proc_detailed:
		m.model_detailed.Refresh()
	default:
		m.updateTable(nil)
	}
}

func (m *model_processes) navDown() {
	switch m.subpage {
	case subpage_proc_main:
		pid, err := m.getSelectedPid()
		if err != nil {
			return // Can't transition
		}
		m.model_detailed.Pid = pid
		m.subpage = subpage_proc_detailed
		m.model_detailed.Refresh()
	}
}

func (m *model_processes) navUp() {
	switch m.subpage {
	case subpage_proc_detailed:
		m.subpage = subpage_proc_main
	}
}

func (m *model_processes) getSelectedPid() (int32, error) {
	row := m.table_processes.SelectedRow()
	if row == nil {
		return 0, errors.New("No processes")
	}
	pid, err := strconv.ParseInt(row[0], 10, 64)
	return int32(pid), err
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
