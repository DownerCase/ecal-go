package publisher

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"

	"github.com/DownerCase/ecal-go/ecal/publisher"
	"github.com/DownerCase/ecal-go/internal/protobuf"
)

// Type must be a pointer and implement the proto.Message interface.
type Msg[T any] interface {
	*T
	proto.Message
}

type Publisher[T proto.Message] struct {
	publisher.Publisher
}

func New[U any, T Msg[U]](topic string) (*Publisher[T], error) {
	var msg T

	pub, err := publisher.New(topic,
		publisher.DataType{
			Name:       protobuf.GetFullName(msg),
			Encoding:   "proto",
			Descriptor: protobuf.GetProtoMessageDescription(msg),
		},
	)
	if err != nil {
		err = fmt.Errorf("protobuf Publisher[%v].New(): %w", reflect.TypeFor[T](), err)
	}

	return &Publisher[T]{*pub}, err
}

func (p *Publisher[T]) Send(t T) error {
	msg, err := proto.Marshal(t)
	if err != nil {
		return fmt.Errorf("protobuf Publisher[%v].Send(): %w", reflect.TypeFor[T](), err)
	}

	p.Messages <- msg

	return nil
}
