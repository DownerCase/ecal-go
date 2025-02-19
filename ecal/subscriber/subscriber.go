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

	"github.com/DownerCase/ecal-go/ecal/ecaltypes"
)

var (
	ErrFailedNew  = errors.New("failed to create new subscriber")
	ErrRcvTimeout = errors.New("timed out")
	ErrRcvBadType = errors.New("receive could not handle type")
)

type GenericSubscriber[T any] struct {
	ecaltypes.Subscriber
	Messages    chan T
	handle      cgo.Handle
	stopped     bool
	Deserialize func(unsafe.Pointer, int) T
}

func NewGenericSubscriber[T any](
	topic string,
	datatype ecaltypes.DataType,
	deserializer func(unsafe.Pointer, int) T,
) (*GenericSubscriber[T], error) {
	sub := &GenericSubscriber[T]{
		Messages:    make(chan T),
		stopped:     false,
		Deserialize: deserializer,
	}
	sub.Subscriber.Callback = sub.subCallback
	handle := cgo.NewHandle(sub.Subscriber)
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

func (sub *GenericSubscriber[T]) Delete() {
	if !bool(C.DestroySubscriber(C.uintptr_t(sub.handle))) {
		// "Failed to delete subscriber"
		return
	}

	if !sub.stopped {
		sub.stopped = true
		close(sub.Messages)
		sub.Messages = nil
	}
}

// Receive a new message from the eCAL receive callback.
func (sub *GenericSubscriber[T]) Receive(timeout time.Duration) (T, error) {
	select {
	case msg := <-sub.Messages:
		return msg, nil
	case <-time.After(timeout):
		var t T
		return t, fmt.Errorf("[Receive]: %w", ErrRcvTimeout)
	}
}

func (sub *GenericSubscriber[T]) subCallback(data unsafe.Pointer, dataLen int) {
	// We must deserialize _before_ submitting the message otherwise
	// the channel may be closed before we finish deserializing
	msg := sub.Deserialize(data, dataLen)
	select {
	case sub.Messages <- msg:
	default:
	}
}

// This function is called by the C code whenever a new message is received
// and deserializes it into a []byte
// If the subscriber Receive is not waiting the incoming message will be dropped
// C.GoBytes takes an int as its length.
//
//export goReceiveCallback
func goReceiveCallback(handle C.uintptr_t, data unsafe.Pointer, dataLen C.int) {
	h := cgo.Handle(handle)
	subBase := h.Value().(ecaltypes.Subscriber)
	subBase.Callback(data, int(dataLen))
}
