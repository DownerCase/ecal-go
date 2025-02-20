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
	"sync/atomic"

	"github.com/DownerCase/ecal-go/ecal/ecaltypes"
)

var ErrFailedNew = errors.New("failed to create new publisher")

type GenericPublisher[T any] struct {
	Messages   chan T
	handle     unsafe.Pointer
	stopped    atomic.Bool
	closed     chan bool
	Serializer func(T) []byte
}

func NewGenericPublisher[T any](
	topic string,
	datatype ecaltypes.DataType,
	serializer func(T) []byte,
) (*GenericPublisher[T], error) {
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

	pub := &GenericPublisher[T]{
		Messages:   make(chan T),
		closed:     make(chan bool),
		handle:     handle,
		Serializer: serializer,
	}

	go pub.sendMessages()

	return pub, nil
}

func (p *GenericPublisher[T]) Delete() {
	if p.stopped.CompareAndSwap(false,true) {
		close(p.Messages)
		<-p.closed // Wait for sendMessages to finish
	}

	C.DestroyPublisher(p.handle)

	// Deleted, clear handle
	p.handle = nil
}

func (p *GenericPublisher[T]) IsStopped() bool {
	return p.stopped.Load()
}

func (p *GenericPublisher[T]) Send(msg T) {
	p.Messages <- msg
}

func (p *GenericPublisher[T]) sendMessages() {
	for msg := range p.Messages {
		bytesMsg := p.Serializer(msg)
		// #cgo noescape PublisherSend
		// #cgo nocallback PublisherSend
		C.PublisherSend(p.handle, unsafe.Pointer(&bytesMsg[0]), C.size_t(len(bytesMsg)))
	}
	p.closed <- true
}
