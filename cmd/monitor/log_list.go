package main

import (
	"time"

	"github.com/DownerCase/ecal-go/ecal/logging"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type LoggingPage int

const (
	subpage_logging_main LoggingPage = iota
	subpage_logging_detailed
)

type logsKeyMap struct {
	table.KeyMap
	Clear key.Binding
}

type model_logs struct {
	table_logs table.Model
	subpage    LoggingPage
	help       help.Model
	keymap     logsKeyMap
	// model_detailed
}

func newLogsKeyMap() logsKeyMap {
	return logsKeyMap{
		KeyMap: table.DefaultKeyMap(),
		Clear: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "Clear"),
		),
	}
}

func (km logsKeyMap) ShortHelp() []key.Binding {
	return append(km.KeyMap.ShortHelp(), km.Clear)
}

func (km logsKeyMap) FullHelp() [][]key.Binding {
	return append([][]key.Binding{{km.Clear}}, km.KeyMap.FullHelp()...)
}

func NewLogsModel() model_logs {
	columns := []table.Column{
		{Title: "Time", Width: 10},
		{Title: "Level", Width: 6},
		{Title: "Unit", Width: 15},
		{Title: "Message", Width: 46},
	}

	t := table.New(
		table.WithHeight(8),
		table.WithFocused(true),
		table.WithColumns(columns),
		table.WithStyles(tableStyle),
	)

	return model_logs{
		table_logs: t,
		subpage:    subpage_logging_main,
		help:       help.New(),
		keymap:     newLogsKeyMap(),
	}
}

func (m *model_logs) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m.subpage {
	case subpage_logging_main:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keymap.Clear):
				m.table_logs.SetRows([]table.Row{})
				m.updateTable(nil)
			default:
				m.updateTable(msg)
			}
		}
	case subpage_logging_detailed:
		// cmd = m.model_detailed.Update(msg)
	}
	return cmd
}

func (m *model_logs) View() string {
	switch m.subpage {
	case subpage_logging_main:
		return baseStyle.Render(m.table_logs.View()) + "\n" + m.help.View(m.keymap)
	case subpage_logging_detailed:
		// return m.model_detailed.View()
	}
	return "Invalid page"
}

func (m *model_logs) Refresh() {
	switch m.subpage {
	case subpage_logging_detailed:
		// m.model_detailed.Refresh()
	default:
		m.updateTable(nil)
	}
}

func (m *model_logs) updateTable(msg tea.Msg) {
	rows := []table.Row{}
	logs := logging.GetLogging().Messages

	for _, log := range logs {
		rows = append(rows, logToRow(log))
	}
	m.table_logs.SetRows(append(m.table_logs.Rows(), rows...))
	m.table_logs, _ = m.table_logs.Update(msg)
}

func logToRow(log logging.LogMessage) table.Row {
	return []string{
		time.UnixMicro(log.Time).Format(time.TimeOnly),
		log.Level.String(),
		log.Unit_name,
		log.Content,
	}
}
