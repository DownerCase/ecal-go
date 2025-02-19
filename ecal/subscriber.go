package ecal

import "C"

import (
	"fmt"
	"reflect"
	"unsafe"

	"google.golang.org/protobuf/proto"

	"github.com/DownerCase/ecal-go/ecal/ecaltypes"
	"github.com/DownerCase/ecal-go/ecal/subscriber"
	"github.com/DownerCase/ecal-go/internal/protobuf"
)

type BinarySubscriber = subscriber.GenericSubscriber[[]byte]

func NewBinarySubscriber(topic string) (*BinarySubscriber, error) {
	sub, err := subscriber.NewGenericSubscriber(topic,
		ecaltypes.DataType{},
		func(data unsafe.Pointer, dataLen int) []byte {
			return C.GoBytes(data, C.int(dataLen))
		},
	)
	if err != nil {
		err = fmt.Errorf("NewBinarySubscriber(): %w", err)
	}

	return sub, err
}

type StringSubscriber = subscriber.GenericSubscriber[string]

func NewStringSubscriber(topic string) (*StringSubscriber, error) {
	sub, err := subscriber.NewGenericSubscriber(topic,
		ecaltypes.DataType{
			Name:     "std::string",
			Encoding: "base",
		},
		func(data unsafe.Pointer, dataLen int) string {
			return C.GoStringN((*C.char)(data), C.int(dataLen))
		},
	)
	if err != nil {
		err = fmt.Errorf("string Subscriber.New(): %w", err)
	}

	return sub, err
}

// Type must be a pointer and implement the proto.Message interface.
type ProtoMessage[T any] interface {
	*T
	proto.Message
}

type ProtobufSubscriber[T proto.Message] = subscriber.GenericSubscriber[T]

func NewProtobufSubscriber[U any, T ProtoMessage[U]](topic string) (*ProtobufSubscriber[T], error) {
	var msg T

	sub, err := subscriber.NewGenericSubscriber(topic,
		ecaltypes.DataType{
			Name:       protobuf.GetFullName(msg),
			Encoding:   "proto",
			Descriptor: protobuf.GetProtoMessageDescription(msg),
		},
		protoDeserialize[U, T],
	)
	if err != nil {
		err = fmt.Errorf("protobuf Subscriber[%v].New(): %w", reflect.TypeFor[U](), err)
	}

	return sub, err
}

func protoDeserialize[U any, T ProtoMessage[U]](data unsafe.Pointer, dataLen int) T {
	// WARNING: Creates a Go slice backed by C data and deserializes into a Go
	// value which gets put into the channel
	bytesUnsafe := unsafe.Slice((*byte)(data), dataLen)

	var msg U

	err := proto.Unmarshal(bytesUnsafe, T(&msg))
	if err != nil {
		// TODO: Don't panic
		panic(fmt.Errorf("protobuf Subscriber[%v].deserialize(): %w", reflect.TypeFor[T](), err))
	}

	return T(&msg)
}
