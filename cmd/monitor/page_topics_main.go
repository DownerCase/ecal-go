package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
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

type ModelTopicsMain struct {
	table  table.Model
	keymap topicsKeyMap
	help   help.Model
	filter entityFilter
}

func NewTopicsMainModel() *ModelTopicsMain {
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

	return &ModelTopicsMain{
		table:  NewTable(cols),
		keymap: newTopicsKeyMap(),
		help:   help.New(),
	}
}

func (m *ModelTopicsMain) Refresh() {
	m.updateTopicsTable(nil)
}

func (m *ModelTopicsMain) Update(msg tea.Msg) tea.Cmd {
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

func (m *ModelTopicsMain) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.help.View(m.keymap)
}

func (m *ModelTopicsMain) GetSelectedID() (uint64, TopicType, error) {
	row := m.table.SelectedRow()
	if row == nil {
		return 0, 0, errEmptyTable
	}

	var topicType TopicType

	switch row[1] {
	case "S":
		topicType = topicTypeSubscriber
	case "P":
		topicType = topicTypePublisher
	}

	id, err := strconv.ParseUint(row[0], 10, 64)
	if err != nil {
		err = fmt.Errorf("topics - GetSelectedID(): %w", err)
	}

	return id, topicType, err
}

func (m *ModelTopicsMain) updateTopicsTable(msg tea.Msg) {
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
		strconv.FormatUint(topic.TopicID, 10),
		strings.ToUpper(topic.Direction[0:1]),
		topic.TopicName,
		topic.Datatype.Name,
		strconv.FormatInt(topic.DataClock, 10),
		strconv.FormatInt(int64(topic.TopicSize), 10),
		strconv.FormatFloat(float64(topic.DataFreq)/1000, 'f', -1, 32),
		strconv.FormatInt(int64(topic.RegistrationClock), 10),
	}
}
