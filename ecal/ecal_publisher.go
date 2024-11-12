package ecal

// #cgo LDFLAGS: -lecal_core
// #include "ecal_go_publisher.h"
import "C"
import "unsafe"
import (
	"fmt"
	"errors"
)

type Publisher struct {
	handle   C.uintptr_t
	Messages chan []byte
	stopped bool
}

func NewPublisher() (Publisher, error) {
	ptr := C.NewPublisher()
	if ptr == nil {
		return Publisher{}, errors.New("Failed to allocate new publisher")
	}
	return Publisher{
		handle:   C.uintptr_t((uintptr(ptr))),
		Messages: make(chan []byte),
	}, nil
}

func DestroyPublisher(publisher *Publisher) bool {
	fmt.Println("Destroying publisher")
	if !publisher.stopped {
		publisher.stopped = true
		close(publisher.Messages)
	}
	return bool(C.DestroyPublisher(publisher.handle))
}

func (p *Publisher) Create(topic string) error {
	if !C.PublisherCreate(p.handle, C.CString(topic)) {
		return errors.New("Failed to Create publisher")
	}
	go p.sendMessages()
	return nil
}

func (p *Publisher) sendMessages() {
	fmt.Println("Starting sending messages")
	for msg := range p.Messages {
		fmt.Println("Publishing message")
		C.PublisherSend(p.handle, unsafe.Pointer(&msg[0]), C.size_t(len(msg)))
	}
	fmt.Println("Finished messages")
}
