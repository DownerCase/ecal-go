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
	if err != nil {
		err = fmt.Errorf("string Publisher.New(): %w", err)
	}
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
	err := p.Publisher.Create(topic,
		publisher.DataType{
			Name:     "std::string",
			Encoding: "base",
		},
	)
	if err != nil {
		err = fmt.Errorf("string Publisher.Create(): %w", err)
	}
	return err
}
