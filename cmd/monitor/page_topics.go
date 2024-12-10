package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TopicsPage int

const (
	subpage_topic_main TopicsPage = iota
	subpage_topic_detailed
	subpage_topic_messages // TODO: Not implemented
)

type model_topics struct {
	subpage TopicsPage
	pages   map[TopicsPage]PageModel
	NavKeys NavKeyMap
}

func NewTopicsModel() *model_topics {
	return (&model_topics{
		subpage: subpage_topic_main,
		pages: map[TopicsPage]PageModel{
			subpage_topic_main:     NewTopicsMainModel(),
			subpage_topic_detailed: NewDetailedModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *model_topics) navDown() {
	switch m.subpage {
	case subpage_topic_main:
		main_model := m.pages[subpage_topic_main].(*model_topics_main)
		topic, is_subscriber, err := main_model.GetSelectedId()
		if err != nil {
			return // Don't' transition
		}
		detailed := m.pages[subpage_topic_detailed].(*model_topic_detailed)
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

func (m *model_topics) Refresh() {
	m.pages[m.subpage].Refresh()
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
	return m.pages[m.subpage].Update(msg)
}

func (m *model_topics) View() string {
	return m.pages[m.subpage].View()
}
