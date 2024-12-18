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

func (s *Subscriber) Receive(timeout time.Duration) (string, error) {
	select {
	case msg := <-s.Messages:
		return msg.(string), nil
	case <-time.After(timeout):
		return "", fmt.Errorf("[Receive]: %w", subscriber.ErrRcvTimeout)
	}
}

func deserialize(data unsafe.Pointer, dataLen int) any {
	return C.GoStringN((*C.char)(data), C.int(dataLen))
}

func (s *Subscriber) Create(topic string) error {
	return s.Subscriber.Create(topic,
		subscriber.DataType{
			Name:     "std::string",
			Encoding: "base",
		},
	)
}
