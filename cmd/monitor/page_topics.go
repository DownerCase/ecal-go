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

type ModelTopics struct {
	subpage TopicsPage
	pages   map[TopicsPage]PageModel
	NavKeys NavKeyMap
}

func NewTopicsModel() *ModelTopics {
	return (&ModelTopics{
		subpage: subpageTopicMain,
		pages: map[TopicsPage]PageModel{
			subpageTopicMain:     NewTopicsMainModel(),
			subpageTopicDetailed: NewDetailedModel(),
			subpageTopicMessages: NewTopicsMessagesModel(),
		},
		NavKeys: make(NavKeyMap),
	}).Init()
}

func (m *ModelTopics) navDown() {
	if m.subpage == subpageTopicMain {
		mainModel := m.pages[subpageTopicMain].(*ModelTopicsMain)
		topic, topicType, err := mainModel.GetSelectedID()
		if err != nil {
			return // Don't' transition
		}
		detailed := m.pages[subpageTopicDetailed].(*ModelTopicDetailed)
		detailed.ShowTopic(topic, topicType)
		m.subpage = subpageTopicDetailed
	}
}

func (m *ModelTopics) navUp() {
	m.subpage = subpageTopicMain
}

func (m *ModelTopics) navMessages() tea.Cmd {
	if m.subpage != subpageTopicMain {
		return nil
	}
	mainModel := m.pages[subpageTopicMain].(*ModelTopicsMain)
	topic, topicType, err := mainModel.GetSelectedID()
	if err != nil {
		return nil // Don't' transition
	}
	messagesModel := m.pages[subpageTopicMessages].(*ModelTopicMessages)
	messagesModel.ShowTopic(topic, topicType)
	m.subpage = subpageTopicMessages
	return messagesModel.Init()
}

func (m *ModelTopics) Refresh() {
	m.pages[m.subpage].Refresh()
}

func (m *ModelTopics) Init() *ModelTopics {
	m.NavKeys["esc"] = func() tea.Cmd { m.navUp(); return nil }
	m.NavKeys["enter"] = func() tea.Cmd { m.navDown(); return nil }
	m.NavKeys["m"] = func() tea.Cmd { return m.navMessages() }
	return m
}

func (m *ModelTopics) Update(msg tea.Msg) tea.Cmd {
	if cmd, navigated := m.NavKeys.HandleMsg(msg); navigated {
		return cmd
	}
	return m.pages[m.subpage].Update(msg)
}

func (m *ModelTopics) View() string {
	return m.pages[m.subpage].View()
}
