package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/DownerCase/ecal-go/ecal/monitoring"
	"github.com/DownerCase/ecal-go/ecal/subscriber"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wrap"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

type ModelTopicMessages struct {
	viewport     viewport.Model
	mon          monitoring.TopicMon
	topicType    TopicType
	topicID      uint64
	subscriber   *subscriber.Subscriber
	msg          []byte
	deserializer func([]byte) string
}

type msgMsg struct {
	msg []byte
}

func NewTopicsMessagesModel() *ModelTopicMessages {
	viewport := viewport.New(85, 9)

	return &ModelTopicMessages{
		viewport: viewport,
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
	default:
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return cmd
}

func (m *ModelTopicMessages) View() string {
	s := strings.Builder{}
	s.WriteString(highlight.Render(m.mon.TopicName))
	s.WriteString(
		fmt.Sprintf(" Messages: %v (%vHz)",
			m.mon.DataClock,
			// TODO: Conversion that can stringify to KHz, MHz, etc.
			strconv.FormatFloat(float64(m.mon.DataFreq)/1000.0, 'e', 3, 64)),
	)
	s.WriteRune('\n')
	// Manually wrap the string with museli/reflow to inject newlines
	// that we can scroll against
	m.viewport.SetContent(wrap.String(m.deserializer(m.msg), m.viewport.Width))
	s.WriteString(m.viewport.View())

	return baseStyle.Render(s.String())
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
		m.deserializer = deserializeBasicString
	case m.mon.Datatype.Encoding == "proto":
		m.deserializer, err = makeProtobufDeserializer(m.mon.Datatype)
		if err != nil {
			panic(err)
		}
	default:
		m.deserializer = deserializeAsHex
	}

	m.msg = nil
	m.subscriber = subscriber
}

func (m *ModelTopicMessages) receiveTicks() tea.Cmd {
	return func() tea.Msg {
		if msg, ok := (<-m.subscriber.Messages).([]byte); ok {
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

func makeProtobufDeserializer(datatype ecal.DataType) (func(msg []byte) string, error) {
	// 1. Take descriptor from the wire and unmarshal it into a descriptorpb
	var descriptorSet descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(datatype.Descriptor, &descriptorSet); err != nil {
		return nil, fmt.Errorf(
			"makeProtobufDeserializer: Failed to unmarshal datatype descriptor %w", err,
		)
	}

	// 2. Turn the file descriptor set into a protoregistry.Files
	registry, err := protodesc.NewFiles(&descriptorSet)
	if err != nil {
		return nil, fmt.Errorf(
			"makeProtobufDeserializer: Failed to parse file descriptor set %w", err,
		)
	}

	// 3. Extract the type descriptors from the registry
	types := dynamicpb.NewTypes(registry)

	// 4. Find the type corresponding to topic datatype
	messageType, err := types.FindMessageByName(protoreflect.FullName(datatype.Name))
	if err != nil {
		return nil, fmt.Errorf(
			"makeProtobufDeserializer: Failed to find definition for message type %s %w",
			protoreflect.FullName(datatype.Name),
			err,
		)
	}

	// 5. Get the message descriptor for the type
	messageDescriptor := messageType.Descriptor()

	// 6. Use the message descriptor to create an instance of a message
	protoMsg := dynamicpb.NewMessage(messageDescriptor)

	// Create a deserializer closure using the extracted message prototype
	return func(msg []byte) string {
		// Unmarshal the binary data into the proto message
		err := proto.Unmarshal(msg, protoMsg)
		if err != nil {
			return fmt.Errorf("protobuf deserialize: Failed to unmarshal %w", err).Error()
		}

		// Marshal the parsed message into a textual format
		textMsg, err := prototext.MarshalOptions{
			Multiline: true,
		}.Marshal(protoMsg)
		if err != nil {
			return fmt.Errorf("protobuf deserialize: Failed to marshal to prototext %w", err).Error()
		}

		// Display that string
		return string(textMsg)
	}, nil
}
