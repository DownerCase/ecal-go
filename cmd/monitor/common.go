package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
)

var (
	errNoTopic    = errors.New("no topic")
	errEmptyTable = errors.New("table empty")
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
	if msg, ok := msg.(tea.KeyMsg); ok {
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

func getTopicFromID(topicType topicType, id string) (monitoring.TopicMon, error) {
	topicList := getTopicMonitoring(topicType)
	for _, topic := range topicList {
		if topic.TopicID == id {
			return topic, nil
		}
	}
	return monitoring.TopicMon{}, fmt.Errorf("[getTopicFromId]: %w", errNoTopic)
}
