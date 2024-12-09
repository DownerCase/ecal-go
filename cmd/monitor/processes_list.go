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
	pages           map[ProcessesPage]PageModel
}

func NewProcessesModel() *model_processes {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "Name", Width: 33},
		{Title: "State", Width: 8},
		{Title: "Info", Width: 23},
		{Title: "Tick", Width: 4},
	}

	pages := make(map[ProcessesPage]PageModel)
	pages[subpage_proc_detailed] = NewDetailedProcessModel()
	return &model_processes{
		table_processes: NewTable(columns),
		subpage:         subpage_proc_main,
		pages:           pages,
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

	if m.subpage == subpage_proc_main {
		return m.updateTable(msg)
	} else {
		return m.pages[m.subpage].Update(msg)
	}
}

func (m *model_processes) View() string {
	if m.subpage == subpage_proc_main {
		return baseStyle.Render(m.table_processes.View()) + "\n" + m.table_processes.HelpView()
	} else {
		return m.pages[m.subpage].View()
	}
}

func (m *model_processes) Refresh() {
	if m.subpage == subpage_proc_main {
		m.updateTable(nil)
	} else {
		m.pages[m.subpage].Refresh()
	}
}

func (m *model_processes) navDown() {
	switch m.subpage {
	case subpage_proc_main:
		pid, err := m.getSelectedPid()
		if err != nil {
			return // Can't transition
		}
		detailed := m.pages[subpage_proc_detailed].(*model_process_detailed)
		detailed.Pid = pid
		m.subpage = subpage_proc_detailed
		detailed.Refresh()
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
