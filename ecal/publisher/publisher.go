package publisher

// #cgo LDFLAGS: -lecal_core
// #cgo CPPFLAGS: -I${SRCDIR}/../../
//#include "publisher.h"
//bool GoPublisherCreate(
//  uintptr_t handle,
//  _GoString_ topic,
//  _GoString_ name, _GoString_ encoding,
//  const char* const descriptor, size_t descriptor_len
//);
import "C"
import "unsafe"
import (
	"errors"

	"github.com/DownerCase/ecal-go/ecal/msg"
)

type DataType = msg.DataType

type Publisher struct {
	Messages chan []byte
	handle   C.uintptr_t
	stopped  bool
}

func New() (*Publisher, error) {
	ptr := C.NewPublisher()
	if ptr == nil {
		return nil, errors.New("Failed to allocate new publisher")
	}
	return &Publisher{
		handle:   C.uintptr_t((uintptr(ptr))),
		Messages: make(chan []byte),
	}, nil
}

func (p *Publisher) Delete() {
	if !p.stopped {
		p.stopped = true
		close(p.Messages)
	}
	if !bool(C.DestroyPublisher(p.handle)) {
		// "Failed to delete publisher"
		return
	}
	// Deleted, clear handle
	p.handle = 0
}

func (p *Publisher) Create(topic string, datatype DataType) error {
	var descriptor_ptr *C.char = nil
	if len(datatype.Descriptor) > 0 {
		descriptor_ptr = (*C.char)(unsafe.Pointer(&datatype.Descriptor[0]))
	}
	if !C.GoPublisherCreate(
		p.handle,
		topic,
		datatype.Name,
		datatype.Encoding,
		descriptor_ptr,
		C.size_t(len(datatype.Descriptor)),
	) {
		return errors.New("Failed to Create publisher")
	}
	go p.sendMessages()
	return nil
}

func (p *Publisher) IsStopped() bool {
	return p.stopped
}

func (p *Publisher) sendMessages() {
	for msg := range p.Messages {
		C.PublisherSend(p.handle, unsafe.Pointer(&msg[0]), C.size_t(len(msg)))
	}
}
