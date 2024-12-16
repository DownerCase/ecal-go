package main

import (
	"errors"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func NewTable(columns []table.Column) table.Model {
	return table.New(
		table.WithHeight(8),
		table.WithFocused(true),
		table.WithStyles(tableStyle),
		table.WithColumns(columns),
	)
}

type NavKeyMap map[string]func() tea.Cmd

func (navKeys NavKeyMap) HandleMsg(msg tea.Msg) (cmd tea.Cmd, navigated bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if f, ok := navKeys[msg.String()]; ok {
			return f(), true
		}
	}
	return nil, false
}

type topicType int

const (
	topicTypeSubscriber topicType = iota
	topicTypePublisher
)

func getTopicMonitoring(topicType topicType) []monitoring.TopicMon {
	switch topicType {
	case topicTypeSubscriber:
		return monitoring.GetMonitoring(monitoring.MonitorSubscriber).Subscribers
	case topicTypePublisher:
		return monitoring.GetMonitoring(monitoring.MonitorPublisher).Publishers
	}
	return nil
}

func getTopicFromId(topicType topicType, id string) (monitoring.TopicMon, error) {
	topic_list := getTopicMonitoring(topicType)
	for _, topic := range topic_list {
		if topic.Topic_id == id {
			return topic, nil
		}
	}
	return monitoring.TopicMon{}, errors.New("Unable to find topic")
}
