package subscriber

import "C"

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"google.golang.org/protobuf/proto"

	"github.com/DownerCase/ecal-go/ecal/subscriber"
	"github.com/DownerCase/ecal-go/internal/protobuf"
)

// Type must be a pointer and implement the proto.Message interface.
type Msg[T any] interface {
	*T
	proto.Message
}

// Both the concrete type and its proto.Message implementing pointer version
// are required to be able to both deserialize and create new values to return.
type Subscriber[U any, T Msg[U]] struct {
	subscriber.Subscriber
}

func New[U any, T Msg[U]](topic string) (*Subscriber[U, T], error) {
	var msg T

	sub, err := subscriber.New(topic,
		subscriber.DataType{
			Name:       protobuf.GetFullName(msg),
			Encoding:   "proto",
			Descriptor: protobuf.GetProtoMessageDescription(msg),
		},
	)
	if err != nil {
		err = fmt.Errorf("protobuf Subscriber[%v].New(): %w", reflect.TypeFor[T](), err)
	}

	sub.Deserialize = deserialize[U, T]
	psub := &Subscriber[U, T]{*sub}

	return psub, err
}

func (s *Subscriber[U, T]) Receive(timeout time.Duration) (U, error) {
	var u U

	var msg any
	select {
	case msg = <-s.Messages:
	case <-time.After(timeout):
		return u, fmt.Errorf("[Receive[%v]()]: %w", reflect.TypeFor[U](), subscriber.ErrRcvTimeout)
	}

	switch msg := msg.(type) {
	case error:
		return u, msg
	case U:
		return msg, nil
	default:
		return u, fmt.Errorf("%w: %v", subscriber.ErrRcvBadType, reflect.TypeOf(msg))
	}
}

func deserialize[U any, T Msg[U]](data unsafe.Pointer, dataLen int) any {
	// WARNING: Creates a Go slice backed by C data and deserializes into a Go
	// value which gets put into the channel
	bytesUnsafe := unsafe.Slice((*byte)(data), dataLen)

	var msg U

	err := proto.Unmarshal(bytesUnsafe, T(&msg))
	if err != nil {
		return fmt.Errorf("protobuf Subscriber[%v].deserialize(): %w", reflect.TypeFor[T](), err)
	}

	return msg
}
