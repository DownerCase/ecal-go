package publisher

// #cgo LDFLAGS: -lecal_core
// #include "publisher.h"
//	bool GoPublisherCreate(
//		uintptr_t handle,
//		_GoString_ topic,
//		_GoString_ name,
//		_GoString_ encoding,
//		const char* const descriptor,
//		size_t descriptor_len
//	) {
//		return PublisherCreate(
//			handle,
//			_GoStringPtr(topic), _GoStringLen(topic),
//			_GoStringPtr(name), _GoStringLen(name),
//			_GoStringPtr(encoding), _GoStringLen(encoding),
//			descriptor, descriptor_len
//		);
//	}
import "C"
import "unsafe"
import (
	"errors"
	"fmt"
)

type Publisher struct {
	Messages chan []byte
	handle   C.uintptr_t
	stopped  bool
}

type DataType struct {
	Name       string
	Encoding   string
	Descriptor []byte
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

func NewStringDataType() DataType {
	return DataType{
		Name:     "std::string",
		Encoding: "base",
	}
}

func (p *Publisher) Delete() {
	fmt.Println("Deleting publisher")
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

func (p *Publisher) sendMessages() {
	for msg := range p.Messages {
		C.PublisherSend(p.handle, unsafe.Pointer(&msg[0]), C.size_t(len(msg)))
	}
}
