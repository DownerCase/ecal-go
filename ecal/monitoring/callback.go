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
			HostName:             C.GoString(pub.host_name),
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

func copyToMethodMons(cmethods []C.struct_CMethodMon) []MethodMon {
	methods := make([]MethodMon, len(cmethods))
	for idx, cmethod := range cmethods {
		methods[idx] = MethodMon{
			Name: C.GoString(cmethod.name),
			RequestType: methodType{
				Type:       C.GoString(cmethod.request_name),
				Descriptor: C.GoString(cmethod.request_desc),
			},
			ResponseType: methodType{
				Type:       C.GoString(cmethod.response_name),
				Descriptor: C.GoString(cmethod.response_desc),
			},
			CallCount: int64(cmethod.call_count),
		}
	}
	return methods
}

func copyToServiceBase(cbase C.struct_CServiceCommon) ServiceBase {
	return ServiceBase{
		Name:              C.GoString(cbase.name),
		Id:                C.GoString(cbase.id),
		RegistrationClock: int32(cbase.registration_clock),
		HostName:          C.GoString(cbase.host_name),
		Process:           C.GoString(cbase.process_name),
		Unit:              C.GoString(cbase.unit_name),
		Pid:               int32(cbase.pid),
		ProtocolVersion:   uint32(cbase.protocol_version),
		Methods:           copyToMethodMons(unsafe.Slice(cbase.methods, cbase.methods_len)),
	}
}

func copyToServerMons(cservers []C.struct_CServerMon) (servers []ServerMon) {
	servers = make([]ServerMon, len(cservers))
	for idx, cserver := range cservers {
		servers[idx] = ServerMon{
			ServiceBase: copyToServiceBase(cserver.base),
			PortV0:      uint32(cserver.port_v0),
			PortV1:      uint32(cserver.port_v1),
		}
	}
	return
}

func copyToClientMons(cclients []C.struct_CClientMon) (clients []ClientMon) {
	clients = make([]ClientMon, len(cclients))
	for idx, cclient := range cclients {
		clients[idx] = ClientMon{
			ServiceBase: copyToServiceBase(cclient.base),
		}
	}
	return
}

//export goCopyMonitoring
func goCopyMonitoring(handle C.uintptr_t, cmon *C.struct_CMonitoring) {
	m := cgo.Handle(handle).Value().(*Monitoring)

	m.Publishers = copyToTopicMons(unsafe.Slice(cmon.publishers, cmon.publishers_len))
	m.Subscribers = copyToTopicMons(unsafe.Slice(cmon.subscribers, cmon.subscribers_len))
	m.Processes = copyToProcessMons(unsafe.Slice(cmon.processes, cmon.processes_len))
	m.Clients = copyToClientMons(unsafe.Slice(cmon.clients, cmon.clients_len))
	m.Servers = copyToServerMons(unsafe.Slice(cmon.servers, cmon.servers_len))
}
