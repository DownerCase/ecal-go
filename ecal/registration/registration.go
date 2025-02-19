package registration

// #cgo LDFLAGS: -lecal_core
//#include "registration.h"
//#include "cgo_types.h"
// #cgo CPPFLAGS: -I${SRCDIR}/..
import "C"

import (
	"log"
	"runtime/cgo"

	"github.com/DownerCase/ecal-go/ecal/ecaltypes"
)

type Event uint8

const (
	EntityNew     Event = 0
	EntityDeleted Event = 1
)

type CallbackToken struct {
	ecalToken uint
	goHandle  cgo.Handle
}

func AddPublisherEventCallback(callback func(ecaltypes.TopicID, Event)) CallbackToken {
	handle := cgo.NewHandle(callback)
	ecalToken := C.AddPublisherEventCallback(C.uintptr_t(handle))
	token := CallbackToken{
		ecalToken: uint(ecalToken),
		goHandle:  handle,
	}

	return token
}

func RemPublisherCallback(token CallbackToken) {
	C.RemPublisherEventCallback(C.uintptr_t(token.ecalToken))
	token.goHandle.Delete()
}

func AddSubscriberEventCallback(callback func(ecaltypes.TopicID, Event)) CallbackToken {
	handle := cgo.NewHandle(callback)
	ecalToken := C.AddSubscriberEventCallback(C.uintptr_t(handle))
	token := CallbackToken{
		ecalToken: uint(ecalToken),
		goHandle:  handle,
	}

	return token
}

func RemSubscriberCallback(token CallbackToken) {
	C.RemSubscriberEventCallback(C.uintptr_t(token.ecalToken))
	token.goHandle.Delete()
}

func toTopicID(id *C.struct_CTopicId) ecaltypes.TopicID {
	return ecaltypes.TopicID{
		TopicID: ecaltypes.EntityID{
			EntityID:  uint64(id.topic_id.entity_id),
			ProcessID: int32(id.topic_id.process_id),
			HostName:  C.GoString(id.topic_id.host_name),
		},
		TopicName: C.GoString(id.topic_name),
	}
}

//export goTopicEventCallback
func goTopicEventCallback(handle C.uintptr_t, id C.struct_CTopicId, event C.uint8_t) {
	h := cgo.Handle(handle)

	f, ok := h.Value().(func(ecaltypes.TopicID, Event))
	if !ok {
		log.Panic("Invalid handle passed to registration callback")
	}

	f(toTopicID(&id), Event(event))
}
