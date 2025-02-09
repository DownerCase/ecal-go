package monitoring

//#include "monitoring.h"
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal"
)

func copyToDatatype(datatype C.struct_CDatatype) ecal.DataType {
	return ecal.DataType{
		Name:       C.GoString(datatype.name),
		Encoding:   C.GoString(datatype.encoding),
		Descriptor: C.GoBytes(datatype.descriptor, datatype.descriptor_len),
	}
}

func copyToTopicMons(ctopics []C.struct_CTopicMon) []TopicMon {
	topics := make([]TopicMon, len(ctopics))
	for idx, pub := range ctopics {
		topics[idx] = TopicMon{
			TopicID:             uint64(pub.topic_id),
			RegistrationClock:   int32(pub.registration_clock),
			TopicName:           C.GoString(pub.topic_name),
			DataClock:           int64(pub.data_clock),
			DataFreq:            int32(pub.data_freq),
			TopicSize:           int32(pub.topic_size),
			UnitName:            C.GoString(pub.unit_name),
			Direction:           C.GoString(pub.direction),
			Datatype:            copyToDatatype(pub.datatype),
			ConnectionsLocal:    int32(pub.connections_local),
			ConnectionsExternal: int32(pub.connections_external),
			MessageDrops:        int32(pub.message_drops),
			HostName:            C.GoString(pub.host_name),
		}
	}

	return topics
}

func copyToProcessMons(cprocs []C.struct_CProcessMon) []ProcessMon {
	procs := make([]ProcessMon, len(cprocs))
	for idx, proc := range cprocs {
		procs[idx] = ProcessMon{
			RegistrationClock:     int32(proc.registration_clock),
			HostName:              C.GoString(proc.host_name),
			Pid:                   int32(proc.pid),
			ProcessName:           C.GoString(proc.process_name),
			UnitName:              C.GoString(proc.unit_name),
			ProcessParameters:     C.GoString(proc.process_parameters), // Command line args
			StateSeverity:         ecal.ProcessSeverity(proc.state_severity),
			StateSeverityLevel:    ecal.ProcessSeverityLevel(proc.state_severity_level),
			StateInfo:             C.GoString(proc.state_info),
			ComponentsInitialized: C.GoString(proc.components),
			RuntimeVersion:        C.GoString(proc.runtime), // eCAL Version in use
		}
	}

	return procs
}

func copyToMethodMons(cmethods []C.struct_CMethodMon) []MethodMon {
	methods := make([]MethodMon, len(cmethods))
	for idx, cmethod := range cmethods {
		methods[idx] = MethodMon{
			Name:         C.GoString(cmethod.name),
			RequestType:  copyToDatatype(cmethod.req_datatype),
			ResponseType: copyToDatatype(cmethod.resp_datatype),
			CallCount:    int64(cmethod.call_count),
		}
	}

	return methods
}

func copyToServiceBase(cbase C.struct_CServiceCommon) ServiceBase {
	return ServiceBase{
		Name:              C.GoString(cbase.name),
		ID:                uint64(cbase.id),
		RegistrationClock: int32(cbase.registration_clock),
		HostName:          C.GoString(cbase.host_name),
		Process:           C.GoString(cbase.process_name),
		Unit:              C.GoString(cbase.unit_name),
		Pid:               int32(cbase.pid),
		ProtocolVersion:   uint32(cbase.protocol_version),
		Methods:           copyToMethodMons(unsafe.Slice(cbase.methods, cbase.methods_len)),
	}
}

func copyToServerMons(cservers []C.struct_CServerMon) []ServerMon {
	servers := make([]ServerMon, len(cservers))
	for idx, cserver := range cservers {
		servers[idx] = ServerMon{
			ServiceBase: copyToServiceBase(cserver.base),
			PortV0:      uint32(cserver.port_v0),
			PortV1:      uint32(cserver.port_v1),
		}
	}

	return servers
}

func copyToClientMons(cclients []C.struct_CClientMon) []ClientMon {
	clients := make([]ClientMon, len(cclients))
	for idx, cclient := range cclients {
		clients[idx] = ClientMon{
			ServiceBase: copyToServiceBase(cclient.base),
		}
	}

	return clients
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
