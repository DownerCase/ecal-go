package main

import (
	"strconv"
	"strings"

	"github.com/DownerCase/ecal-go/ecal/monitoring"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type topicsKeyMap struct {
	table.KeyMap
	FilterAll key.Binding
	FilterPub key.Binding
	FilterSub key.Binding
	Messages  key.Binding
	Help      key.Binding
}

type entityFilter uint8

const (
	entityAll        entityFilter = 0
	entityPublisher  entityFilter = 1
	entitySubscriber entityFilter = 2
)

func newTopicsKeyMap() topicsKeyMap {
	return topicsKeyMap{
		KeyMap: table.DefaultKeyMap(),
		FilterAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "All"),
		),
		FilterPub: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Publishers"),
		),
		FilterSub: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "Subscribers"),
		),
		Messages: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "Messages"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Help"),
		),
	}
}

func (km topicsKeyMap) ShortHelp() []key.Binding {
	return append(km.KeyMap.ShortHelp(), km.FilterAll, km.FilterPub, km.FilterSub, km.Messages)
}

func (km topicsKeyMap) FullHelp() [][]key.Binding {
	return append([][]key.Binding{{km.FilterAll, km.FilterPub, km.FilterSub}}, km.KeyMap.FullHelp()...)
}

type modelTopicsMain struct {
	table  table.Model
	keymap topicsKeyMap
	help   help.Model
	filter entityFilter
}

func NewTopicsMainModel() *modelTopicsMain {
	cols := []table.Column{
		{Title: "ID", Width: 0}, // Zero width ID column to use as identifier
		{Title: "D", Width: 1},
		{Title: "Topic", Width: 20},
		{Title: "Type", Width: 20},
		{Title: "Count", Width: 8},
		{Title: "Size", Width: 8},
		{Title: "Freq (Hz)", Width: 10},
		{Title: "Tick", Width: 4},
	}

	return &modelTopicsMain{
		table:  NewTable(cols),
		keymap: newTopicsKeyMap(),
		help:   help.New(),
	}
}

func (m *modelTopicsMain) Refresh() {
	m.updateTopicsTable(nil)
}

func (m *modelTopicsMain) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keymap.FilterAll):
			m.filter = entityAll
			m.updateTopicsTable(nil)
		case key.Matches(msg, m.keymap.FilterPub):
			m.filter = entityPublisher
			m.updateTopicsTable(nil)
		case key.Matches(msg, m.keymap.FilterSub):
			m.filter = entitySubscriber
			m.updateTopicsTable(nil)
		default:
			m.updateTopicsTable(msg)
		}
	default:
		m.updateTopicsTable(msg)
	}
	return nil
}

func (m *modelTopicsMain) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.help.View(m.keymap)
}

func (m *modelTopicsMain) GetSelectedID() (string, topicType, error) {
	row := m.table.SelectedRow()
	if row == nil {
		return "", 0, errEmptyTable
	}
	var topicType topicType
	switch row[1] {
	case "S":
		topicType = topicTypeSubscriber
	case "P":
		topicType = topicTypePublisher
	}
	return row[0], topicType, nil
}

func (m *modelTopicsMain) updateTopicsTable(msg tea.Msg) {
	rows := []table.Row{}
	entities := monitoring.MonitorNone
	switch m.filter {
	case entityAll:
		entities = monitoring.MonitorPublisher | monitoring.MonitorSubscriber
	case entitySubscriber:
		entities = monitoring.MonitorSubscriber
	case entityPublisher:
		entities = monitoring.MonitorPublisher
	}
	mon := monitoring.GetMonitoring(entities)
	for _, topic := range mon.Publishers {
		rows = append(rows, topicToRow(topic))
	}
	for _, topic := range mon.Subscribers {
		rows = append(rows, topicToRow(topic))
	}
	m.table.SetRows(rows)
	m.table, _ = m.table.Update(msg)
}

func topicToRow(topic monitoring.TopicMon) table.Row {
	return []string{
		topic.TopicID,
		strings.ToUpper(topic.Direction[0:1]),
		topic.TopicName,
		topic.Datatype.Name,
		strconv.FormatInt(topic.DataClock, 10),
		strconv.FormatInt(int64(topic.TopicSize), 10),
		strconv.FormatFloat(float64(topic.DataFreq)/1000, 'f', -1, 32),
		strconv.FormatInt(int64(topic.RegistrationClock), 10),
	}
}
