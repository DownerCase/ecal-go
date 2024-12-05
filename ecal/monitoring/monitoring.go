package monitoring

//#cgo LDFLAGS: -lecal_core
//#include <ecal/ecal_process_severity.h>
//#include "monitoring.h"
//#cgo CPPFLAGS: -I${SRCDIR}/../types
import "C"
import (
	"runtime/cgo"
	"strconv"

	"github.com/DownerCase/ecal-go/ecal/types"
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

type ProcessSeverity uint8

const (
	ProcSevUnknown  ProcessSeverity = C.proc_sev_unknown
	ProcSevHealthy  ProcessSeverity = C.proc_sev_healthy
	ProcSevWarning  ProcessSeverity = C.proc_sev_warning
	ProcSevCritical ProcessSeverity = C.proc_sev_critical
	ProcSevFailed   ProcessSeverity = C.proc_sev_failed
)

func (p ProcessSeverity) String() string {
	switch p {
	case ProcSevUnknown:
		return "Unknown"
	case ProcSevHealthy:
		return "Healthy"
	case ProcSevWarning:
		return "Warning"
	case ProcSevCritical:
		return "Critical"
	case ProcSevFailed:
		return "Failed"
	default:
		return strconv.FormatUint(uint64(p), 10)
	}
}

type TopicMon struct {
	Registration_clock int32 // registration heart beat
	// host_name          string
	// host_group         string
	// pid                int32
	// process_name       string
	Unit_name  string
	Topic_id   string
	Topic_name string
	Direction  string
	Datatype   types.DataType
	// TODO: transport layer
	Topic_size           int32 // Size of messages (bytes)
	Connections_local    int32
	Connections_external int32
	Message_drops        int32
	// data_id              int64
	Data_clock int64
	Data_freq  int32 // mHz
	// attributes
}

type ProcessMon struct {
	Registration_clock   int32 // registration heart beat
	Host_name            string
	Host_group           string
	Pid                  int32
	Process_name         string
	Unit_name            string
	Process_parameters   string // Command line args
	State_severity       ProcessSeverity
	State_severity_level int32
	State_info           string
	// TODO: Time sync?
	Components_initialized string
	Runtime_version        string // eCAL Version in use
}

type methodType struct {
	Type       string
	Descriptor string
}

type MethodMon struct {
	Name         string
	RequestType  methodType
	ResponseType methodType
	CallCount    int64
}

type ServiceBase struct {
	Name              string
	Id                string
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
