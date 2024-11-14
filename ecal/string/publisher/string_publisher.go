package publisher

import (
	"github.com/DownerCase/ecal-go/ecal/publisher"
)

type Publisher struct {
	publisher.Publisher
}

func New() (*Publisher, error) {
	pub, err := publisher.New()
	return &Publisher{*pub}, err
}

func (p *Publisher) Send(msg string) error {
	p.Messages <- []byte(msg)
	return nil
}

func (p *Publisher) Create(topic string) error {
	return p.Publisher.Create(topic,
		publisher.DataType{
			Name:     "std::string",
			Encoding: "base",
		},
	)
}
