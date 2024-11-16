package subscriber

// #cgo LDFLAGS: -lecal_core
// #cgo CPPFLAGS: -I${SRCDIR}/../../
// #include "subscriber.h"
//	bool GoSubscriberCreate(
//		uintptr_t handle,
//		_GoString_ topic,
//		_GoString_ name, _GoString_ encoding,
//		const char* const descriptor, size_t descriptor_len
//	);
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal/msg"
)

type Subscriber struct {
	messages chan []byte
	handle   cgo.Handle
	stopped  bool
}

type DataType = msg.DataType

func New() (*Subscriber, error) {
	sub := Subscriber{
		messages: make(chan []byte),
		stopped:  false,
	}
	handle := cgo.NewHandle(sub)
	sub.handle = handle
	if !C.NewSubscriber(C.uintptr_t(sub.handle)) {
		handle.Delete()
		return nil, errors.New("Failed to allocate new subscriber")
	}
	return &sub, nil
}

func (p *Subscriber) Delete() {
	if !p.stopped {
		p.stopped = true
		close(p.messages)
	}
	if !bool(C.DestroySubscriber(C.uintptr_t(p.handle))) {
		// "Failed to delete subscriber"
		return
	}
	// Deleted, clear handle
	p.handle.Delete()
}

func (p *Subscriber) Create(topic string, datatype DataType) error {
	var descriptor_ptr *C.char = nil
	if len(datatype.Descriptor) > 0 {
		descriptor_ptr = (*C.char)(unsafe.Pointer(&datatype.Descriptor[0]))
	}
	if !C.GoSubscriberCreate(
		C.uintptr_t(p.handle),
		topic,
		datatype.Name,
		datatype.Encoding,
		descriptor_ptr,
		C.size_t(len(datatype.Descriptor)),
	) {
		return errors.New("Failed to Create publisher")
	}
	return nil
}

// Receive a new message from the eCAL receive callback
func (p *Subscriber) Receive() []byte {
	return <-p.messages
}

// Deserialize straight from the eCAL internal buffer to our Go []byte
func (p *Subscriber) deserialize(data unsafe.Pointer, len C.long) []byte {
	return C.GoBytes(data, C.int(len))
}

// This function is called by the C code whenever a new message is received
// and deserializes it into a []byte
// If the subscriber Receive is not waiting the incoming message will be dropped
//
//export goReceiveCallback
func goReceiveCallback(handle C.uintptr_t, data unsafe.Pointer, len C.long) {
	h := cgo.Handle(handle)
	sub := h.Value().(Subscriber)
	select {
	case sub.messages <- sub.deserialize(data, len):
	default:
	}
}
