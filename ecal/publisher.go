package ecal

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"

	"github.com/DownerCase/ecal-go/ecal/ecaltypes"
	"github.com/DownerCase/ecal-go/ecal/publisher"
	"github.com/DownerCase/ecal-go/internal/protobuf"
)

type BinaryPublisher = publisher.GenericPublisher[[]byte]

func NewBinaryPublisher(topic string) (*BinaryPublisher, error) {
	pub, err := publisher.NewGenericPublisher(
		topic,
		ecaltypes.DataType{},
		func(msg []byte) []byte { return msg },
	)
	if err != nil {
		err = fmt.Errorf("NewBinaryPublisher(): %w", err)
	}

	return pub, err
}

type StringPublisher = publisher.GenericPublisher[string]

func NewStringPublisher(topic string) (*StringPublisher, error) {
	pub, err := publisher.NewGenericPublisher(
		topic,
		ecaltypes.DataType{
			Name:     "std::string",
			Encoding: "base",
		},
		func(msg string) []byte { return []byte(msg) },
	)
	if err != nil {
		err = fmt.Errorf("NewStringPublisher(): %w", err)
	}

	return pub, err
}

type ProtobufPublisher[T proto.Message] = publisher.GenericPublisher[T]

func NewProtobufPublisher[U any, T ProtoMessage[U]](topic string) (*ProtobufPublisher[T], error) {
	var msg T

	pub, err := publisher.NewGenericPublisher(topic,
		ecaltypes.DataType{
			Name:       protobuf.GetFullName(msg),
			Encoding:   "proto",
			Descriptor: protobuf.GetProtoMessageDescription(msg),
		},
		func(msg T) []byte {
			protoMsg, err := proto.Marshal(msg)
			if err != nil {
				panic(fmt.Errorf("protobuf Publisher[%v].Send(): %w", reflect.TypeFor[T](), err))
			}

			return protoMsg
		},
	)
	if err != nil {
		err = fmt.Errorf("protobuf Publisher[%v].New(): %w", reflect.TypeFor[T](), err)
	}

	return pub, err
}
