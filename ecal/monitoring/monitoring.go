package monitoring

//#cgo LDFLAGS: -lecal_core
//#include "monitoring.h"
//#cgo CPPFLAGS: -I${SRCDIR}/../types
import "C"
import (
	"runtime/cgo"

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

type Monitoring struct {
	Publishers  []TopicMon
	Subscribers []TopicMon
}

type Callback = func()

func GetMonitoring(entities MonitorEntity) Monitoring {
	var mon Monitoring
	handle := cgo.NewHandle(&mon)
	C.GetMonitoring(C.uintptr_t(handle), C.uint(entities))
	handle.Delete()
	return mon
}
