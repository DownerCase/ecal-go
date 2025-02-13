package publisher

// #cgo LDFLAGS: -lecal_core
// #include "publisher.h"
// void *GoNewPublisher(
//  _GoString_ topic,
//  _GoString_ name, _GoString_ encoding,
//  const char* const descriptor, size_t descriptor_len
// );
// // C preamble.
import "C"

import (
	"errors"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal"
)

var ErrFailedNew = errors.New("failed to create new publisher")

type DataType = ecal.DataType

type Publisher struct {
	Messages chan []byte
	handle   unsafe.Pointer
	stopped  bool
	closed   chan bool
}

func New(topic string, datatype DataType) (*Publisher, error) {
	var descriptorPtr *C.char
	if len(datatype.Descriptor) > 0 {
		descriptorPtr = (*C.char)(unsafe.Pointer(&datatype.Descriptor[0]))
	}

	handle := C.GoNewPublisher(
		topic,
		datatype.Name,
		datatype.Encoding,
		descriptorPtr,
		C.size_t(len(datatype.Descriptor)),
	)
	if handle == nil {
		return nil, ErrFailedNew
	}

	pub := &Publisher{
		Messages: make(chan []byte),
		stopped:  false,
		closed:   make(chan bool),
		handle:   handle,
	}

	go pub.sendMessages()

	return pub, nil
}

func (p *Publisher) Delete() {
	if !p.stopped {
		p.stopped = true
		close(p.Messages)
		<-p.closed // Wait for sendMessages to finish
	}

	C.DestroyPublisher(p.handle)

	// Deleted, clear handle
	p.handle = nil
}

func (p *Publisher) IsStopped() bool {
	return p.stopped
}

func (p *Publisher) sendMessages() {
	for msg := range p.Messages {
		C.PublisherSend(p.handle, unsafe.Pointer(&msg[0]), C.size_t(len(msg)))
	}
	p.closed <- true
}
