package registration

// #cgo LDFLAGS: -lecal_core
//#include "registration.h"
//#include "types.h"
// #cgo CPPFLAGS: -I${SRCDIR}/../types
import "C"

import (
	"runtime/cgo"

	"github.com/DownerCase/ecal-go/ecal/types"
)

type Event uint8

const (
	ENTITY_NEW     Event = 0
	ENTITY_DELETED Event = 1
)

type QualityFlags uint8

type (
	EntityId = types.EntityId
	TopicId  = types.TopicId
)

type CallbackToken struct {
	ecal_token uint
	go_handle  cgo.Handle
}

type QualityTopicInfo struct {
	Datatype     types.DataType
	QualityFlags QualityFlags
}

func AddPublisherEventCallback(callback func(TopicId, Event)) CallbackToken {
	handle := cgo.NewHandle(callback)
	ecal_token := C.AddPublisherEventCallback(C.uintptr_t(handle))
	token := CallbackToken{
		ecal_token: uint(ecal_token),
		go_handle:  handle,
	}
	return token
}

func RemPublisherCallback(token CallbackToken) {
	C.RemPublisherEventCallback(C.uintptr_t(token.ecal_token))
	token.go_handle.Delete()
}

func AddSubscriberEventCallback(callback func(TopicId, Event)) CallbackToken {
	handle := cgo.NewHandle(callback)
	ecal_token := C.AddSubscriberEventCallback(C.uintptr_t(handle))
	token := CallbackToken{
		ecal_token: uint(ecal_token),
		go_handle:  handle,
	}
	return token
}

func RemSubscriberCallback(token CallbackToken) {
	C.RemSubscriberEventCallback(C.uintptr_t(token.ecal_token))
	token.go_handle.Delete()
}

func toTopicId(id *C.struct_CTopicId) TopicId {
	return TopicId{
		Topic_id: EntityId{
			Entity_id:  C.GoString(id.topic_id.entity_id),
			Process_id: int32(id.topic_id.process_id),
			Host_name:  C.GoString(id.topic_id.host_name),
		},
		Topic_name: C.GoString(id.topic_name),
	}
}

//export goTopicEventCallback
func goTopicEventCallback(handle C.uintptr_t, id C.struct_CTopicId, event C.uint8_t) {
	h := cgo.Handle(handle)
	f := h.Value().(func(TopicId, Event))
	f(toTopicId(&id), Event(event))
}
