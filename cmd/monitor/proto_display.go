package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/DownerCase/ecal-go/cmd/monitor/internal/protobuf/protorange"
	"github.com/DownerCase/ecal-go/ecal"
)

func makeProtobufDeserializer(datatype ecal.DataType) (func(msg []byte) string, error) {
	m, err := createBlankMessage(datatype)

	return func(msg []byte) string {
		// Unmarshal the binary data into the proto message
		err := proto.Unmarshal(msg, m)
		if err != nil {
			return fmt.Errorf("protobuf deserialize: Failed to unmarshal %w", err).Error()
		}

		s := strings.Builder{}

		err = protorange.Options{
			Stable:            true,
			EmitDefaultValues: true,
		}.Range(m.ProtoReflect(),
			writePush(&s),
			writePop(&s),
		)
		if err != nil {
			panic(err)
		}

		return s.String()
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

func writePush(s *strings.Builder) func(protopath.Values) error {
	return func(p protopath.Values) error {
		// Print the key.
		var fd protoreflect.FieldDescriptor

		beforeLast := p.Index(-2)

		last := p.Index(-1)
		indent := strings.Repeat("\t", max(0, p.Len()-2))

		switch last.Step.Kind() {
		case protopath.FieldAccessStep:
			fd = last.Step.FieldDescriptor()
			s.WriteString(fmt.Sprintf("%s%s: ", indent, fd.Name()))
		case protopath.ListIndexStep:
			fd = beforeLast.Step.FieldDescriptor() // lists always appear in the context of a repeated field
			s.WriteString(fmt.Sprintf("%s%d: ", indent, last.Step.ListIndex()))
		case protopath.MapIndexStep:
			fd = beforeLast.Step.FieldDescriptor() // maps always appear in the context of a repeated field
			s.WriteString(fmt.Sprintf("%s%v: ", indent, last.Step.MapIndex().Interface()))
		case protopath.AnyExpandStep:
			s.WriteString(fmt.Sprintf("%s[%v]: ", indent, last.Value.Message().Descriptor().FullName()))
		case protopath.UnknownAccessStep:
			s.WriteString(indent + "?: ")
		case protopath.RootStep:
			// Don't indent the root node
			return nil
		}

		// Starting printing the value.
		switch v := last.Value.Interface().(type) {
		case protoreflect.Message:
			s.WriteString("{\n")
		case protoreflect.List:
			s.WriteString("[\n")
		case protoreflect.Map:
			s.WriteString("{\n")
		case protoreflect.EnumNumber:
			var ev protoreflect.EnumValueDescriptor
			if fd != nil {
				ev = fd.Enum().Values().ByNumber(v)
			}

			if ev != nil {
				s.WriteString(fmt.Sprintf("%v\n", ev.Name()))
			} else {
				s.WriteString(fmt.Sprintf("%v\n", v))
			}
		case string, []byte:
			s.WriteString(fmt.Sprintf("%q\n", v))
		default:
			s.WriteString(fmt.Sprintf("%v\n", v))
		}

		return nil
	}
}

func writePop(s *strings.Builder) func(protopath.Values) error {
	return func(p protopath.Values) error {
		// Finish printing the value.
		last := p.Index(-1)
		switch last.Step.Kind() {
		case protopath.RootStep:
			// Don't dedent the root node
			return nil
		case protopath.FieldAccessStep:
		case protopath.ListIndexStep:
		case protopath.MapIndexStep:
		case protopath.AnyExpandStep:
		case protopath.UnknownAccessStep:
		}

		indent := strings.Repeat("\t", p.Len()-2)
		switch last.Value.Interface().(type) {
		case protoreflect.Message:
			s.WriteString(indent + "}\n")
		case protoreflect.List:
			s.WriteString(indent + "]\n")
		case protoreflect.Map:
			s.WriteString(indent + "}\n")
		}

		return nil
	}
}
