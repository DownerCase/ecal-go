package main

import (
	"fmt"
	"strings"

	"github.com/DownerCase/ecal-go/ecal"
	"github.com/charmbracelet/bubbles/table"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/DownerCase/ecal-go/cmd/monitor/internal/protobuf/protorange"
)

type protoReflector struct {
	items       []table.Row
	isCollapsed func(string) bool
}

func makeProtobufDeserializer(
	datatype ecal.DataType,
	isCollapsed func(string) bool) (func(msg []byte,
) []table.Row, error,
) {
	m, err := createBlankMessage(datatype)

	return func(msg []byte) []table.Row {
		// Unmarshal the binary data into the proto message
		err := proto.Unmarshal(msg, m)
		if err != nil {
			return []table.Row{{"", fmt.Errorf("protobuf deserialize: Failed to unmarshal %w", err).Error()}}
		}

		pr := protoReflector{
			items:       make([]table.Row, 0),
			isCollapsed: isCollapsed,
		}

		err = protorange.Options{
			Stable:            true,
			EmitDefaultValues: true,
		}.Range(m.ProtoReflect(),
			pr.writePush,
			pr.writePop,
		)
		if err != nil {
			panic(err)
		}

		return pr.items
	}, err
}

func createBlankMessage(datatype ecal.DataType) (*dynamicpb.Message, error) {
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
	return dynamicpb.NewMessage(messageDescriptor), nil
}

// Yes this function is long...
func (reflector *protoReflector) writePush(p protopath.Values) error { //nolint:funlen
	// Print the key.
	// If the parent is collapsed, break
	if reflector.isCollapsed(p.Path[:len(p.Path)-1].String()) {
		return protorange.Break
	}

	var fd protoreflect.FieldDescriptor

	s := strings.Builder{}
	last := p.Index(-1)
	beforeLast := p.Index(-2)
	indent := strings.Repeat("\t", max(0, p.Len()-2))

	switch last.Step.Kind() {
	case protopath.RootStep:
		// Don't indent the root node
		return nil
	case protopath.FieldAccessStep:
		fd = last.Step.FieldDescriptor()
		s.WriteString(fmt.Sprintf("%s%s: ", indent, fd.Name()))
	case protopath.ListIndexStep:
		// lists always appear in the context of a repeated field
		fd = beforeLast.Step.FieldDescriptor()
		s.WriteString(fmt.Sprintf("%s%d: ", indent, last.Step.ListIndex()))
	case protopath.MapIndexStep:
		// maps always appear in the context of a repeated field
		fd = beforeLast.Step.FieldDescriptor()
		s.WriteString(fmt.Sprintf("%s%v: ", indent, last.Step.MapIndex().Interface()))
	case protopath.AnyExpandStep:
		s.WriteString(fmt.Sprintf("%s[%v]: ", indent, last.Value.Message().Descriptor().FullName()))
	case protopath.UnknownAccessStep:
		s.WriteString(indent + "?: ")
	}

	collapsed := reflector.isCollapsed(p.Path.String())

	// Starting printing the value.
	switch v := last.Value.Interface().(type) {
	case protoreflect.Message:
		s.WriteString("{")

		if collapsed {
			s.WriteString(
				fmt.Sprintf(" %d fields }",
					last.Value.Message().Descriptor().Fields().Len(),
				),
			)
		}
	case protoreflect.List:
		s.WriteString("[")

		if collapsed {
			s.WriteString(
				fmt.Sprintf(" %d items ]",
					last.Value.List().Len(),
				),
			)
		}
	case protoreflect.Map:
		s.WriteString("{")

		if collapsed {
			s.WriteString(
				fmt.Sprintf(" %d entries }",
					last.Value.Message().Descriptor().Fields().Len(),
				),
			)
		}
	case protoreflect.EnumNumber:
		var ev protoreflect.EnumValueDescriptor
		if fd != nil {
			ev = fd.Enum().Values().ByNumber(v)
		}

		if ev != nil {
			s.WriteString(fmt.Sprintf("%v", ev.Name()))
		} else {
			s.WriteString(fmt.Sprintf("%v", v))
		}
	case string, []byte:
		s.WriteString(fmt.Sprintf("%q", v))
	default:
		s.WriteString(fmt.Sprintf("%v", v))
	}

	if lineItem := s.String(); len(lineItem) > 0 {
		reflector.items = append(reflector.items, table.Row{p.Path.String(), lineItem})
	}

	return nil
}

func (reflector *protoReflector) writePop(p protopath.Values) error {
	if reflector.isCollapsed(p.Path.String()) {
		return nil
	}

	if reflector.isCollapsed(p.Path[:len(p.Path)-1].String()) {
		return nil
	}
	// Finish printing the value.
	last := p.Index(-1)
	switch last.Step.Kind() {
	case protopath.RootStep:
		// Don't do anything for the root node
		return nil
	case protopath.FieldAccessStep:
	case protopath.ListIndexStep:
	case protopath.MapIndexStep:
	case protopath.AnyExpandStep:
	case protopath.UnknownAccessStep:
	}

	// -1 for the root node, -1 for the scope we're leaving
	indent := strings.Repeat("\t", p.Len()-2)
	lineItem := ""

	switch last.Value.Interface().(type) {
	case protoreflect.Message:
		lineItem = indent + "}"
	case protoreflect.List:
		lineItem = indent + "]"
	case protoreflect.Map:
		lineItem = indent + "}"
	}

	if len(lineItem) > 0 {
		reflector.items = append(reflector.items, table.Row{p.Path.String(), lineItem})
	}

	return nil
}
