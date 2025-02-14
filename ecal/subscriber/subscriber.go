package subscriber

// #cgo LDFLAGS: -lecal_core
// #cgo CPPFLAGS: -I${SRCDIR}/../../
// #include "subscriber.h"
//	bool GoNewSubscriber(
//		uintptr_t handle,
//		_GoString_ topic,
//		_GoString_ name, _GoString_ encoding,
//		const char* const descriptor, size_t descriptor_len
//	);
import "C"

import (
	"errors"
	"fmt"
	"runtime/cgo"
	"time"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal"
)

var (
	ErrFailedNew  = errors.New("failed to create new subscriber")
	ErrRcvTimeout = errors.New("timed out")
	ErrRcvBadType = errors.New("receive could not handle type")
)

type Subscriber struct {
	Messages    chan any
	handle      cgo.Handle
	stopped     bool
	Deserialize func(unsafe.Pointer, int) any
}

type DataType = ecal.DataType

func New(topic string, datatype DataType) (*Subscriber, error) {
	sub := &Subscriber{
		Messages:    make(chan any),
		stopped:     false,
		Deserialize: deserializer,
	}
	handle := cgo.NewHandle(sub)
	sub.handle = handle

	var descriptorPtr *C.char
	if len(datatype.Descriptor) > 0 {
		descriptorPtr = (*C.char)(unsafe.Pointer(&datatype.Descriptor[0]))
	}

	if !C.GoNewSubscriber(
		C.uintptr_t(sub.handle),
		topic,
		datatype.Name,
		datatype.Encoding,
		descriptorPtr,
		C.size_t(len(datatype.Descriptor)),
	) {
		handle.Delete()
		return nil, ErrFailedNew
	}

	return sub, nil
}

func (p *Subscriber) Delete() {
	if !bool(C.DestroySubscriber(C.uintptr_t(p.handle))) {
		// "Failed to delete subscriber"
		return
	}

	if !p.stopped {
		p.stopped = true
		close(p.Messages)
		p.Messages = nil
	}
}

// Receive a new message from the eCAL receive callback.
func (p *Subscriber) Receive(timeout time.Duration) ([]byte, error) {
	select {
	case msg := <-p.Messages:
		return msg.([]byte), nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("[Receive]: %w", ErrRcvTimeout)
	}
}

// Deserialize straight from the eCAL internal buffer to our Go []byte.
func deserializer(data unsafe.Pointer, dataLen int) any {
	return C.GoBytes(data, C.int(dataLen))
}

// This function is called by the C code whenever a new message is received
// and deserializes it into a []byte
// If the subscriber Receive is not waiting the incoming message will be dropped
// C.GoBytes takes an int as its length
//
//export goReceiveCallback
func goReceiveCallback(handle C.uintptr_t, data unsafe.Pointer, dataLen C.int) {
	h := cgo.Handle(handle)
	sub := h.Value().(*Subscriber)
	// We must deserialize _before_ submitting the message otherwise
	// the channel may be closed before we finish deserializing
	msg := sub.Deserialize(data, int(dataLen))
	select {
	case sub.Messages <- msg:
	default:
	}
}
