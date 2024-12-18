package subscriber

import "C"

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal/subscriber"
)

type Subscriber struct {
	subscriber.Subscriber
}

func New() (*Subscriber, error) {
	sub, err := subscriber.New()
	sub.Deserialize = deserialize
	return &Subscriber{*sub}, err
}

func (p *Subscriber) Receive(timeout time.Duration) (string, error) {
	select {
	case msg := <-p.Messages:
		return msg.(string), nil
	case <-time.After(timeout):
		return "", fmt.Errorf("[Receive]: %w", subscriber.ErrRcvTimeout)
	}
}

func deserialize(data unsafe.Pointer, len int) any {
	return C.GoStringN((*C.char)(data), C.int(len))
}

func (s *Subscriber) Create(topic string) error {
	return s.Subscriber.Create(topic,
		subscriber.DataType{
			Name:     "std::string",
			Encoding: "base",
		},
	)
}
