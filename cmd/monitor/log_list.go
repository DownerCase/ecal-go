package main

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/DownerCase/ecal-go/ecal/logging"
)

type LoggingPage int

const (
	subpageLoggingMain LoggingPage = iota
	subpageLoggingDetailed
)

type logsKeyMap struct {
	table.KeyMap
	Clear key.Binding
}

type modelLogs struct {
	table   table.Model
	subpage LoggingPage
	help    help.Model
	keymap  logsKeyMap
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

func NewLogsModel() *modelLogs {
	columns := []table.Column{
		{Title: "Time", Width: 10},
		{Title: "Level", Width: 6},
		{Title: "Unit", Width: 15},
		{Title: "Message", Width: 46},
	}

	return &modelLogs{
		table:   NewTable(columns),
		subpage: subpageLoggingMain,
		help:    help.New(),
		keymap:  newLogsKeyMap(),
	}
}

func (m *modelLogs) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m.subpage {
	case subpageLoggingMain:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keymap.Clear):
				m.table.SetRows([]table.Row{})
				m.updateTable(nil)
			default:
				m.updateTable(msg)
			}
		}
	case subpageLoggingDetailed:
		// cmd = m.model_detailed.Update(msg)
	}
	return cmd
}

func (m *modelLogs) View() string {
	switch m.subpage {
	case subpageLoggingMain:
		return baseStyle.Render(m.table.View()) + "\n" + m.help.View(m.keymap)
	case subpageLoggingDetailed:
		// return m.model_detailed.View()
	}
	return "Invalid page"
}

func (m *modelLogs) Refresh() {
	switch m.subpage {
	case subpageLoggingDetailed:
		// m.model_detailed.Refresh()
	default:
		m.updateTable(nil)
	}
}

func (m *modelLogs) updateTable(msg tea.Msg) {
	rows := []table.Row{}
	logs := logging.GetLogging().Messages

	for _, log := range logs {
		rows = append(rows, logToRow(log))
	}
	m.table.SetRows(append(m.table.Rows(), rows...))
	m.table, _ = m.table.Update(msg)
}

func logToRow(log logging.LogMessage) table.Row {
	return []string{
		time.UnixMicro(log.Time).Format(time.TimeOnly),
		log.Level.String(),
		log.UnitName,
		log.Content,
	}
}
