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
)

type TopicsPage int

const (
	subpage_topic_main TopicsPage = iota
	subpage_topic_detailed
	subpage_topic_messages // TODO: Not implemented
)

type model_topics struct {
	table_topics table.Model
	keymap       topicsKeyMap
	help         help.Model
	filter       entityFilter
	subpage      TopicsPage
	pages        map[TopicsPage]PageModel
	NavKeys      NavKeyMap
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

func NewTopicsModel() *model_topics {
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

	return (&model_topics{
		table_topics: NewTable(topics_columns),
		keymap:       newTopicsKeyMap(),
		help:         help.New(),
		subpage:      subpage_topic_main,
		pages: map[TopicsPage]PageModel{
			subpage_topic_detailed: NewDetailedModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
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

func (m *model_topics) navDown() {
	switch m.subpage {
	case subpage_topic_main:
		topic, is_subscriber, err := m.GetSelectedId()
		if err != nil {
			return // Don't' transition
		}
		detailed := m.pages[subpage_topic_detailed].(*model_detailed)
		detailed.ShowTopic(topic, is_subscriber)
		m.subpage = subpage_topic_detailed
	}
}

func (m *model_topics) navUp() {
	switch m.subpage {
	case subpage_topic_detailed:
		m.subpage = subpage_topic_main
	}
}

func (m *model_topics) GetSelectedId() (string, bool, error) {
	row := m.table_topics.SelectedRow()
	if row == nil {
		return "", false, errors.New("No active topics")
	}
	return row[0], row[1] == "S", nil
}

func (m *model_topics) Refresh() {
	if m.subpage == subpage_topic_main {
		m.updateTopicsTable(nil)
	} else {
		m.pages[m.subpage].Refresh()
	}
}

func (m *model_topics) Init() *model_topics {
	m.NavKeys[tea.KeyEscape] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys[tea.KeyEnter] = func() tea.Cmd { m.navDown(); return nil }
	return m
}

func (m *model_topics) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}

	var cmd tea.Cmd
	if m.subpage == subpage_topic_main {
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
	} else {
		cmd = m.pages[m.subpage].Update(msg)
	}
	return cmd
}

func (m *model_topics) View() string {
	if m.subpage == subpage_topic_main {
		return baseStyle.Render(m.table_topics.View()) + "\n" + m.help.View(m.keymap)
	} else {
		return m.pages[m.subpage].View()
	}
}
