package publisher

import (
	"fmt"

	"github.com/DownerCase/ecal-go/ecal/publisher"
)

type Publisher struct {
	publisher.Publisher
}

func New() (*Publisher, error) {
	pub, err := publisher.New()
	return &Publisher{*pub}, err
}

// Send a message formatted with fmt.Print.
func (p *Publisher) Send(msg ...any) error {
	p.Messages <- []byte(fmt.Sprint(msg...))
	return nil
}

// Send a message formatted with fmt.Printf.
func (p *Publisher) Sendf(format string, a ...any) error {
	p.Messages <- []byte(fmt.Sprintf(format, a...))
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
