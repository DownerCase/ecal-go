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
import "unsafe"
import (
	"errors"

	"github.com/DownerCase/ecal-go/ecal/msg"
)

type Subscriber struct {
	messages chan []byte
	handle   C.uintptr_t
	stopped  bool
}

type DataType = msg.DataType

func New() (*Subscriber, error) {
	ptr := C.NewSubscriber()
	if ptr == nil {
		return nil, errors.New("Failed to allocate new subscriber")
	}
	return &Subscriber{
		handle:   C.uintptr_t((uintptr(ptr))),
		messages: make(chan []byte),
	}, nil
}

func (p *Subscriber) Delete() {
	if !p.stopped {
		p.stopped = true
		close(p.messages)
	}
	if !bool(C.DestroySubscriber(p.handle)) {
		// "Failed to delete subscriber"
		return
	}
	// Deleted, clear handle
	p.handle = 0
}

func (p *Subscriber) Create(topic string, datatype DataType) error {
	var descriptor_ptr *C.char = nil
	if len(datatype.Descriptor) > 0 {
		descriptor_ptr = (*C.char)(unsafe.Pointer(&datatype.Descriptor[0]))
	}
	if !C.GoSubscriberCreate(
		p.handle,
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

// Receive a new message from eCAL
// Currently performs at least two copies
// - Internal eCAL recieve buffer -> ReceiveBuffer's buffer in C wrapper
// - C wrapper -> C.GoBytes
// TODO: Use a callback based method to copy the data directly from eCAL's
// buffer to a Go []byte result variable
func (p *Subscriber) Receive() []byte {
	var msg *C.char
	var len C.size_t
	// WARNING: Calling through cgo three times in a frequently run function
	// is suboptimal

	// Receive message
	handle := C.SubscriberReceive(p.handle, &msg, &len)
	// Copy to a go []byte
	go_msg := C.GoBytes(unsafe.Pointer(msg), C.int(len))
	// Free the original receive buffer
	C.DestroyMessage(handle)
	return go_msg
}
