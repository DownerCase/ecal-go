package publisher

// #cgo LDFLAGS: -lecal_core
// #cgo CPPFLAGS: -I${SRCDIR}/../../
//#include "publisher.h"
// bool GoNewPublisher(
//  uintptr_t handle,
//  _GoString_ topic,
//  _GoString_ name, _GoString_ encoding,
//  const char* const descriptor, size_t descriptor_len
// );
import "C"

import (
	"errors"
	"runtime/cgo"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal/types"
)

var ErrFailedNew = errors.New("failed to create new publisher")

type DataType = types.DataType

type Publisher struct {
	Messages chan []byte
	handle   cgo.Handle
	stopped  bool
}

func New(topic string, datatype DataType) (*Publisher, error) {
	pub := &Publisher{
		Messages: make(chan []byte),
		stopped:  false,
	}
	handle := cgo.NewHandle(pub)
	pub.handle = handle

	var descriptorPtr *C.char
	if len(datatype.Descriptor) > 0 {
		descriptorPtr = (*C.char)(unsafe.Pointer(&datatype.Descriptor[0]))
	}

	if !C.GoNewPublisher(
		C.uintptr_t(pub.handle),
		topic,
		datatype.Name,
		datatype.Encoding,
		descriptorPtr,
		C.size_t(len(datatype.Descriptor)),
	) {
		handle.Delete()
		return nil, ErrFailedNew
	}
	go pub.sendMessages()
	return pub, nil
}

func (p *Publisher) Delete() {
	if !p.stopped {
		p.stopped = true
		close(p.Messages)
	}
	if !bool(C.DestroyPublisher(C.uintptr_t(p.handle))) {
		// "Failed to delete publisher"
		return
	}
	// Deleted, clear handle
	p.handle = 0
}

func (p *Publisher) IsStopped() bool {
	return p.stopped
}

func (p *Publisher) sendMessages() {
	for msg := range p.Messages {
		C.PublisherSend(C.uintptr_t(p.handle), unsafe.Pointer(&msg[0]), C.size_t(len(msg)))
	}
}
