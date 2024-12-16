package logging

//#include "logging.h"
//#include "types.h"
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

func copyToLogMessages(cmsgs []C.struct_CLogMessage) []LogMessage {
	msgs := make([]LogMessage, len(cmsgs))
	for idx, msg := range cmsgs {
		msgs[idx] = LogMessage{
			Time:      int64(msg.time),
			Host:      C.GoString(msg.host_name),
			Process:   C.GoString(msg.process_name),
			Unit_name: C.GoString(msg.unit_name),
			Content:   C.GoString(msg.content),
			Pid:       int32(msg.pid),
			Level:     Level(msg.level),
		}
	}
	return msgs
}

//export goCopyLogging
func goCopyLogging(handle C.uintptr_t, clogging *C.struct_CLogging) {
	l := cgo.Handle(handle).Value().(*Logging)

	num_messages := clogging.num_messages
	if num_messages > 0 {
		ms := (*[1 << 30]C.struct_CLogMessage)(unsafe.Pointer(clogging.messages))[:num_messages:num_messages]
		l.Messages = copyToLogMessages(ms)
	}
}
