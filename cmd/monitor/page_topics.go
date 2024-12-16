package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TopicsPage int

const (
	subpageTopicMain TopicsPage = iota
	subpageTopicDetailed
	subpageTopicMessages // TODO: Not implemented
)

type modelTopics struct {
	subpage TopicsPage
	pages   map[TopicsPage]PageModel
	NavKeys NavKeyMap
}

func NewTopicsModel() *modelTopics {
	return (&modelTopics{
		subpage: subpageTopicMain,
		pages: map[TopicsPage]PageModel{
			subpageTopicMain:     NewTopicsMainModel(),
			subpageTopicDetailed: NewDetailedModel(),
			subpageTopicMessages: NewTopicsMessagesModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *modelTopics) navDown() {
	switch m.subpage {
	case subpageTopicMain:
		mainModel := m.pages[subpageTopicMain].(*modelTopicsMain)
		topic, topicType, err := mainModel.GetSelectedID()
		if err != nil {
			return // Don't' transition
		}
		detailed := m.pages[subpageTopicDetailed].(*modelTopicDetailed)
		detailed.ShowTopic(topic, topicType)
		m.subpage = subpageTopicDetailed
	}
}

func (m *modelTopics) navUp() {
	switch m.subpage {
	default:
		m.subpage = subpageTopicMain
	}
}

func (m *modelTopics) navMessages() tea.Cmd {
	if m.subpage != subpageTopicMain {
		return nil
	}
	mainModel := m.pages[subpageTopicMain].(*modelTopicsMain)
	topic, topicType, err := mainModel.GetSelectedID()
	if err != nil {
		return nil // Don't' transition
	}
	messagesModel := m.pages[subpageTopicMessages].(*modelTopicMessages)
	messagesModel.ShowTopic(topic, topicType)
	m.subpage = subpageTopicMessages
	return messagesModel.Init()
}

func (m *modelTopics) Refresh() {
	m.pages[m.subpage].Refresh()
}

func (m *modelTopics) Init() *modelTopics {
	m.NavKeys["esc"] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys["enter"] = func() tea.Cmd { m.navDown(); return nil }
	m.NavKeys["m"] = func() tea.Cmd { return m.navMessages() }
	return m
}

func (m *modelTopics) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}
	return m.pages[m.subpage].Update(msg)
}

func (m *modelTopics) View() string {
	return m.pages[m.subpage].View()
}
