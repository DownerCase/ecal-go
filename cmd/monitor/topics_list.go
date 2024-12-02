package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/DownerCase/ecal-go/ecal/monitoring"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model_topics struct {
	table_topics table.Model
	keymap       topicsKeyMap
	help         help.Model
	filter       entityFilter
}

type topicsKeyMap struct {
	table.KeyMap
	FilterAll key.Binding
	FilterPub key.Binding
	FilterSub key.Binding
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
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Help"),
		),
	}
}

func (km topicsKeyMap) ShortHelp() []key.Binding {
	return append(km.KeyMap.ShortHelp(), km.FilterAll, km.FilterPub, km.FilterSub)
}

func (km topicsKeyMap) FullHelp() [][]key.Binding {
	return append([][]key.Binding{{km.FilterAll, km.FilterPub, km.FilterSub}}, km.KeyMap.FullHelp()...)
}

func NewTopicsModel() model_topics {
	topics_columns := []table.Column{
		{Title: "ID", Width: 0}, // Zero width ID column to use as identifier
		{Title: "D", Width: 1},
		{Title: "Topic", Width: 20},
		{Title: "Type", Width: 20},
		{Title: "Count", Width: 8},
		{Title: "Size", Width: 8},
		{Title: "Freq (Hz)", Width: 10},
		{Title: "Tick", Width: 4},
	}
	t := table.New(
		table.WithHeight(8),
		table.WithFocused(true),
		table.WithColumns(topics_columns),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return model_topics{
		table_topics: t,
		keymap:       newTopicsKeyMap(),
		help:         help.New(),
	}
}

func topicToRow(topic monitoring.TopicMon) table.Row {
	return []string{
		topic.Topic_id,
		strings.ToUpper(topic.Direction[0:1]),
		topic.Topic_name,
		topic.Datatype.Name,
		strconv.FormatInt(topic.Data_clock, 10),
		strconv.FormatInt(int64(topic.Topic_size), 10),
		strconv.FormatFloat(float64(topic.Data_freq)/1000, 'f', -1, 32),
		strconv.FormatInt(int64(topic.Registration_clock), 10),
	}
}

func (m *model_topics) updateTopicsTable(msg tea.Msg) {
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
	m.table_topics.SetRows(rows)
	m.table_topics, _ = m.table_topics.Update(msg)
}

func (m *model_topics) GetSelectedId() (string, bool, error) {
	row := m.table_topics.SelectedRow()
	if row == nil {
		return "", false, errors.New("No active topics")
	}
	return row[0], row[1] == "S", nil
}

func (m *model_topics) Refresh() {
	m.updateTopicsTable(nil)
}

func (m *model_topics) Init() tea.Cmd {
	return nil
}

func (m *model_topics) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
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
	}
	return cmd
}

func (m *model_topics) View() string {
	return baseStyle.Render(m.table_topics.View()) + "\n" + m.help.View(m.keymap)
}
