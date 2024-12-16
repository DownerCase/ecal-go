package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type modelTopicDetailed struct {
	table     table.Model
	id        string    `exhaustruct:"optional"`
	topicType topicType `exhaustruct:"optional"`
}

func NewDetailedModel() *modelTopicDetailed {
	cols := []table.Column{
		{Title: "", Width: 14},
		{Title: "", Width: 67},
	}

	return &modelTopicDetailed{
		table: NewTable(cols),
	}
}

func (m *modelTopicDetailed) ShowTopic(topicID string, topicType topicType) {
	m.id = topicID
	m.topicType = topicType
	m.updateDetailedTable(nil)
}

func (m *modelTopicDetailed) Init() tea.Cmd {
	return nil
}

func (m *modelTopicDetailed) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.table, cmd = m.table.Update(msg)
	}
	return cmd
}

func (m *modelTopicDetailed) View() string {
	return baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}

func (m *modelTopicDetailed) Refresh() {
	m.updateDetailedTable(nil)
}

func (m *modelTopicDetailed) updateDetailedTable(msg tea.Msg) {
	t, _ := getTopicFromID(m.topicType, m.id)
	m.table.Columns()[0].Title = t.Direction
	m.table.Columns()[1].Title = t.TopicName
	rows := []table.Row{
		{"Datatype", fmt.Sprintf("(%s) %s", t.Datatype.Encoding, t.Datatype.Name)},
		{"Unit", t.UnitName},
		{
			"Messages",
			fmt.Sprintf("%v (%v dropped)", t.DataClock, t.MessageDrops),
		},
		{"Frequency", strconv.FormatFloat(float64(t.DataFreq)/1000, 'f', -1, 32)},
		{"Message Size", strconv.FormatInt(int64(t.TopicSize), 10)},
		{
			"Connections",
			fmt.Sprintf("%v local, %v external", t.ConnectionsLocal, t.ConnectionsExternal),
		},
		{"Tick", strconv.FormatInt(int64(t.RegistrationClock), 10)},
	}
	m.table.SetRows(rows)
	m.table, _ = m.table.Update(msg)
}
