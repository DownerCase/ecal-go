package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/DownerCase/ecal-go/ecal/subscriber"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wrap"
)

type ModelTopicMessages struct {
	table        table.Model
	mon          monitoring.TopicMon
	topicType    TopicType
	topicID      uint64
	subscriber   *subscriber.Subscriber
	msg          []byte
	deserializer func([]byte) []table.Row
	collapsed    map[string]struct{}
}

type msgMsg struct {
	msg []byte
}

func NewTopicsMessagesModel() *ModelTopicMessages {
	return &ModelTopicMessages{
		table:     NewTable([]table.Column{{Title: "ID", Width: 0}, {Title: "Initializing", Width: 85}}),
		collapsed: make(map[string]struct{}),
	}
}

func (m *ModelTopicMessages) Init() tea.Cmd {
	return m.receiveTicks()
}

func (m *ModelTopicMessages) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case msgMsg:
		m.msg = msg.msg
		cmd = m.receiveTicks()
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			key := m.table.SelectedRow()[0]
			if _, ok := m.collapsed[key]; ok {
				// Already collapsed, expand
				delete(m.collapsed, key)
			} else {
				// Collapse
				m.collapsed[key] = struct{}{}
			}
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}

	return cmd
}

func (m *ModelTopicMessages) View() string {
	s := strings.Builder{}
	s.WriteString(highlight.Render(m.mon.TopicName))
	s.WriteString(
		fmt.Sprintf(" Message: %v (%vHz)",
			m.mon.DataClock,
			// TODO: Conversion that can stringify to KHz, MHz, etc.
			strconv.FormatFloat(float64(m.mon.DataFreq)/1000.0, 'e', 3, 64)),
	)

	m.table.Columns()[1].Title = s.String()

	// Manually wrap the string with museli/reflow to inject newlines
	// that we can scroll against
	m.table.SetRows(m.deserializer(m.msg))

	return baseStyle.Render(m.table.View())
}

func (m *ModelTopicMessages) Refresh() {
	m.mon, _ = getTopicFromID(m.topicType, m.topicID)
}

func (m *ModelTopicMessages) ShowTopic(topicID uint64, topicType TopicType) {
	if m.topicID != topicID {
		m.topicType = topicType
		m.topicID = topicID
		m.mon, _ = getTopicFromID(m.topicType, m.topicID)
		m.createSubscriber()
	}

	m.Refresh()
}

func (m *ModelTopicMessages) createSubscriber() {
	// (re)create subscriber with new topic type
	if m.subscriber != nil {
		m.subscriber.Delete()
	}

	subscriber, err := subscriber.New(m.mon.TopicName, m.mon.Datatype)
	if err != nil {
		subscriber.Delete()
		panic(fmt.Errorf("[Topic Messages]: %w", err))
	}

	switch {
	case m.mon.Datatype.Name == "std::string" && m.mon.Datatype.Encoding == "base":
		m.deserializer = deserializeBasicString(m.table.Width)
	case m.mon.Datatype.Encoding == "proto":
		// Implemented in dedicated file
		m.deserializer, err = makeProtobufDeserializer(
			m.mon.Datatype,
			func(key string) bool { _, ok := m.collapsed[key]; return ok },
		)
		if err != nil {
			panic(err)
		}
	default:
		m.deserializer = deserializeAsHex(m.table.Width)
	}

	m.msg = nil
	m.subscriber = subscriber
}

func (m *ModelTopicMessages) receiveTicks() tea.Cmd {
	return func() tea.Msg {
		// Hard throttle to 250 updates/s
		time.Sleep(time.Millisecond * 4)

		if msg, ok := (<-m.subscriber.Messages).([]byte); ok {
			return msgMsg{msg: msg}
		}

		return nil
	}
}

func wrapAndSplitToItems(content string, width int) []table.Row {
	// Remove trailing blank line
	if content[len(content)-1] == '\n' {
		content = content[:len(content)-1]
	}
	// Manually wrap the string with museli/reflow to inject newlines
	// that we can scroll against
	wrapped := wrap.String(content, width)
	lines := strings.Split(wrapped, "\n")

	items := make([]table.Row, len(lines))
	for idx, line := range lines {
		items[idx] = table.Row{"", line}
	}

	return items
}

// Message deserializers.
func deserializeBasicString(getWidth func() int) func(msg []byte) []table.Row {
	return func(msg []byte) []table.Row {
		return wrapAndSplitToItems(string(msg), getWidth())
	}
}

func deserializeAsHex(getWidth func() int) func(msg []byte) []table.Row {
	return func(msg []byte) []table.Row {
		content := hex.EncodeToString(msg)
		return wrapAndSplitToItems(content, getWidth())
	}
}
