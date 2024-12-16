package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model_topic_detailed struct {
	table_detailed    table.Model
	detailed_topic_id string
	topicType         topicType
}

func NewDetailedModel() *model_topic_detailed {

	cols := []table.Column{
		{Title: "", Width: 14},
		{Title: "", Width: 67},
	}

	return &model_topic_detailed{
		table_detailed:    NewTable(cols),
		detailed_topic_id: "",
	}
}

func (m *model_topic_detailed) ShowTopic(topic_id string, topicType topicType) {
	m.detailed_topic_id = topic_id
	m.topicType = topicType
	m.updateDetailedTable(nil)
}

func (m *model_topic_detailed) Init() tea.Cmd {
	return nil
}

func (m *model_topic_detailed) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.table_detailed, cmd = m.table_detailed.Update(msg)
	}
	return cmd
}

func (m *model_topic_detailed) View() string {
	return baseStyle.Render(m.table_detailed.View()) + "\n" + m.table_detailed.HelpView()
}

func (m *model_topic_detailed) Refresh() {
	m.updateDetailedTable(nil)
}

func (m *model_topic_detailed) updateDetailedTable(msg tea.Msg) {
	t, _ := getTopicFromId(m.topicType, m.detailed_topic_id)
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
