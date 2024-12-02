package main

import (
	"fmt"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/monitoring"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model_detailed struct {
	table_detailed    table.Model
	detailed_topic_id string
	is_subscriber     bool
}

func NewDetailedModel() model_detailed {

	cols := []table.Column{
		{Title: "", Width: 14},
		{Title: "", Width: 67},
	}

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

	t := table.New(
		table.WithHeight(8),
		table.WithFocused(true),
		table.WithColumns(cols),
		table.WithStyles(s),
	)

	return model_detailed{
		table_detailed:    t,
		detailed_topic_id: "",
	}
}

func (m *model_detailed) ShowTopic(topic_id string, is_subscriber bool) {
	m.detailed_topic_id = topic_id
	m.is_subscriber = is_subscriber
	m.updateDetailedTable(nil)
}

func (m *model_detailed) Init() tea.Cmd {
	return nil
}

func (m *model_detailed) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TickMsg:
		m.updateDetailedTable(nil)
		return doTick()
	case tea.KeyMsg:
		m.table_detailed, cmd = m.table_detailed.Update(msg)
	}
	return cmd
}

func (m *model_detailed) View() string {
	return baseStyle.Render(m.table_detailed.View()) + "\n" + m.table_detailed.HelpView()
}

func (m *model_detailed) updateDetailedTable(msg tea.Msg) {
	mon := monitoring.GetMonitoring(monitoring.MonitorPublisher | monitoring.MonitorSubscriber)
	var t monitoring.TopicMon
	var topic_list []monitoring.TopicMon
	if m.is_subscriber {
		topic_list = mon.Subscribers
	} else {
		topic_list = mon.Publishers
	}
	for _, topic := range topic_list {
		if topic.Topic_id == m.detailed_topic_id {
			t = topic
			break
		}
	}
	m.table_detailed.Columns()[0].Title = t.Direction
	m.table_detailed.Columns()[1].Title = t.Topic_name
	rows := []table.Row{
		{"Datatype", fmt.Sprintf("(%s) %s", t.Datatype.Encoding, t.Datatype.Name)},
		{"Unit", t.Unit_name},
		{"Messages",
			fmt.Sprintf("%v (%v dropped)", t.Data_clock, t.Message_drops),
		},
		{"Frequency", strconv.FormatFloat(float64(t.Data_freq)/1000, 'f', -1, 32)},
		{"Message Size", strconv.FormatInt(int64(t.Topic_size), 10)},
		{"Connections",
			fmt.Sprintf("%v local, %v external", t.Connections_local, t.Connections_external),
		},
		{"Tick", strconv.FormatInt(int64(t.Registration_clock), 10)},
	}
	m.table_detailed.SetRows(rows)
	m.table_detailed, _ = m.table_detailed.Update(msg)
}
