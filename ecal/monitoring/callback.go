package monitoring

//#include "monitoring.h"
import "C"
import (
	"runtime/cgo"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal/types"
)

func copyToTopicMons(ctopics []C.struct_CTopicMon) []TopicMon {
	topics := make([]TopicMon, len(ctopics))
	for idx, pub := range ctopics {
		topics[idx] = TopicMon{
			Topic_id:           C.GoString(pub.topic_id),
			Registration_clock: int32(pub.registration_clock),
			Topic_name:         C.GoString(pub.topic_name),
			Data_clock:         int64(pub.data_clock),
			Data_freq:          int32(pub.data_freq),
			Topic_size:         int32(pub.topic_size),
			Unit_name:          C.GoString(pub.unit_name),
			Direction:          C.GoString(pub.direction),
			Datatype: types.DataType{
				Name:     C.GoString(pub.datatype.name),
				Encoding: C.GoString(pub.datatype.encoding),
			},
			Connections_local:    int32(pub.connections_local),
			Connections_external: int32(pub.connections_external),
			Message_drops:        int32(pub.message_drops),
		}
	}
	return topics
}

func copyToProcessMons(cprocs []C.struct_CProcessMon) []ProcessMon {
	procs := make([]ProcessMon, len(cprocs))
	for idx, proc := range cprocs {
		procs[idx] = ProcessMon{
			Registration_clock:     int32(proc.registration_clock),
			Host_name:              C.GoString(proc.host_name),
			Pid:                    int32(proc.pid),
			Process_name:           C.GoString(proc.process_name),
			Unit_name:              C.GoString(proc.unit_name),
			Process_parameters:     C.GoString(proc.process_parameters), // Command line args
			State_severity:         ProcessSeverity(proc.state_severity),
			State_severity_level:   int32(proc.state_severity_level),
			State_info:             C.GoString(proc.state_info),
			Components_initialized: C.GoString(proc.components),
			Runtime_version:        C.GoString(proc.runtime), // eCAL Version in use
		}
	}
	return procs
}

//export goCopyMonitoring
func goCopyMonitoring(handle C.uintptr_t, cmon *C.struct_CMonitoring) {
	m := cgo.Handle(handle).Value().(*Monitoring)

	numPublishers := cmon.publishers_len
	if numPublishers > 0 {
		p := (*[1 << 30]C.struct_CTopicMon)(unsafe.Pointer(cmon.publishers))[:numPublishers:numPublishers]
		m.Publishers = copyToTopicMons(p)
	}
	numSubscribers := cmon.subscribers_len
	if numSubscribers > 0 {
		s := (*[1 << 30]C.struct_CTopicMon)(unsafe.Pointer(cmon.subscribers))[:numSubscribers:numSubscribers]
		m.Subscribers = copyToTopicMons(s)
	}
	numProcesses := cmon.processes_len
	if numProcesses > 0 {
		p := (*[1 << 30]C.struct_CProcessMon)(unsafe.Pointer(cmon.processes))[:numProcesses:numProcesses]
		m.Processes = copyToProcessMons(p)
	}

}
