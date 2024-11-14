package publisher

import (
	"github.com/DownerCase/ecal-go/ecal/publisher"
	"github.com/DownerCase/ecal-go/internal/protobuf"

	"google.golang.org/protobuf/proto"
)

type Publisher[T proto.Message] struct {
	publisher.Publisher
}

func New[T proto.Message](t T) (*Publisher[T], error) {
	pub, err := publisher.New()
	return &Publisher[T]{*pub}, err
}

func (p *Publisher[T]) Send(t T) error {
	msg, err := proto.Marshal(t)
	if err != nil {
		return err
	}

	p.Messages <- msg
	return nil
}

func (p *Publisher[T]) Create(topic string) error {
	var msg T
	return p.Publisher.Create(topic,
		publisher.DataType{
			Name:       protobuf.GetFullName(msg),
			Encoding:   "proto",
			Descriptor: protobuf.GetProtoMessageDescription(msg),
		},
	)
}
