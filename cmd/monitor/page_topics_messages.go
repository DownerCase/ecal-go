package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/DownerCase/ecal-go/ecal/subscriber"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wrap"
)

type model_topic_messages struct {
	viewport     viewport.Model
	mon          monitoring.TopicMon
	topicType    topicType
	topicId      string
	subscriber   *subscriber.Subscriber
	msg          []byte
	deserializer func([]byte) string
}

type msgMsg struct {
	msg []byte
}

func NewTopicsMessagesModel() *model_topic_messages {
	viewport := viewport.New(85, 10)
	subscriber, _ := subscriber.New()
	return &model_topic_messages{
		viewport:   viewport,
		subscriber: subscriber,
	}
}

func (m *model_topic_messages) Init() tea.Cmd {
	return m.receiveTicks()
}

func (m *model_topic_messages) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case msgMsg:
		m.msg = msg.msg
		cmd = m.receiveTicks()
	default:
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return cmd
}

func (m *model_topic_messages) View() string {
	s := strings.Builder{}
	s.WriteString(highlight.Render(m.mon.Topic_name))
	s.WriteString(
		fmt.Sprintf(" Messages: %v (%vHz)",
			m.mon.Data_clock,
			// TODO: Conversion that can stringify to KHz, MHz, etc.
			strconv.FormatFloat(float64(m.mon.Data_freq)/1000.0, 'e', 3, 64)),
	)
	s.WriteRune('\n')
	// Manually wrap the string with museli/reflow to inject newlines
	// that we can scroll against
	m.viewport.SetContent(wrap.String(m.deserializer(m.msg), m.viewport.Width))
	s.WriteString(m.viewport.View())
	return baseStyle.Render(s.String())
}

func (m *model_topic_messages) Refresh() {
	m.mon, _ = getTopicFromId(m.topicType, m.topicId)
}

func (m *model_topic_messages) ShowTopic(topicId string, topicType topicType) {
	if m.topicId != topicId {
		m.topicType = topicType
		m.topicId = topicId
		m.mon, _ = getTopicFromId(m.topicType, m.topicId)
		m.createSubscriber()
	}
	m.Refresh()
}

func (m *model_topic_messages) createSubscriber() {
	// (re)create subscriber with new topic type
	m.subscriber.Delete()
	subscriber, err := subscriber.New()
	if err != nil {
		subscriber.Delete()
		panic(fmt.Errorf("[Topic Messages]: %w", err))
	}
	err = subscriber.Create(m.mon.Topic_name, m.mon.Datatype)
	if err != nil {
		subscriber.Delete()
		panic(fmt.Errorf("[Topic Messages]: %w", err))
	}
	switch {
	case m.mon.Datatype.Name == "std::string" && m.mon.Datatype.Encoding == "base":
		m.deserializer = deserializeBasicString
	default:
		m.deserializer = deserializeAsHex
	}
	m.msg = nil
	m.subscriber = subscriber
}

func (m *model_topic_messages) receiveTicks() tea.Cmd {
	return func() tea.Msg {
		switch msg := (<-m.subscriber.Messages).(type) {
		case []byte:
			return msgMsg{msg: msg}
		}
		return nil
	}
}

// Message deserializers.
func deserializeBasicString(msg []byte) string {
	return string(msg)
}

func deserializeAsHex(msg []byte) string {
	return hex.EncodeToString(msg)
}
