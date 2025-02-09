package monitoring

//#cgo LDFLAGS: -lecal_core
//#include "monitoring.h"
//#cgo CPPFLAGS: -I${SRCDIR}/..
import "C"

import (
	"runtime/cgo"

	"github.com/DownerCase/ecal-go/ecal"
)

type MonitorEntity uint

const (
	MonitorNone       MonitorEntity = C.monitoring_none
	MonitorPublisher  MonitorEntity = C.monitoring_publisher
	MonitorSubscriber MonitorEntity = C.monitoring_subscriber
	MonitorServer     MonitorEntity = C.monitoring_server
	MonitorClient     MonitorEntity = C.monitoring_client
	MonitorProcess    MonitorEntity = C.monitoring_process
	MonitorHost       MonitorEntity = C.monitoring_host
	MonitorAll        MonitorEntity = C.monitoring_all
)

type TopicMon struct {
	RegistrationClock int32 // registration heart beat
	HostName          string
	// host_group         string
	// pid                int32
	// process_name       string
	UnitName  string
	TopicID   uint64
	TopicName string
	Direction string
	Datatype  ecal.DataType
	// TODO: transport layer
	TopicSize           int32 // Size of messages (bytes)
	ConnectionsLocal    int32
	ConnectionsExternal int32
	MessageDrops        int32
	// data_id              int64
	DataClock int64
	DataFreq  int32 // mHz
	// attributes
}

type ProcessMon struct {
	RegistrationClock  int32 // registration heart beat
	HostName           string
	ShmDomain          string
	Pid                int32
	ProcessName        string
	UnitName           string
	ProcessParameters  string // Command line args
	StateSeverity      ecal.ProcessSeverity
	StateSeverityLevel ecal.ProcessSeverityLevel
	StateInfo          string
	// TODO: Time sync?
	ComponentsInitialized string
	RuntimeVersion        string // eCAL Version in use
}

type MethodMon struct {
	Name         string
	RequestType  ecal.DataType
	ResponseType ecal.DataType
	CallCount    int64
}

type ServiceBase struct {
	Name              string
	ID                uint64
	Methods           []MethodMon
	RegistrationClock int32 // registration heart beat
	HostName          string
	Process           string
	Unit              string
	Pid               int32
	ProtocolVersion   uint32
}

type ServerMon struct {
	ServiceBase
	PortV0 uint32 // TCP Port for V0 protocol
	PortV1 uint32 // TCP Port for V1 protocol
}

type ClientMon struct {
	ServiceBase
}

type Monitoring struct {
	Publishers  []TopicMon
	Subscribers []TopicMon
	Processes   []ProcessMon
	Clients     []ClientMon
	Servers     []ServerMon
}

func GetMonitoring(entities MonitorEntity) Monitoring {
	var mon Monitoring
	handle := cgo.NewHandle(&mon)
	// The C code calls goCopyMonitoring to fill the above Monitoring variable
	// via the handle
	C.GetMonitoring(C.uintptr_t(handle), C.uint(entities))
	handle.Delete()

	return mon
}
