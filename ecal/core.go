package ecal

// #cgo LDFLAGS: -lecal_core
// #include "core.h"
// #include "cgo_config.h"
// #include <stdlib.h>
import "C"

import (
	"strconv"
	"unsafe"

	"github.com/DownerCase/ecal-go/ecal/ecaltypes"
)

type DataType = ecaltypes.DataType

const (
	// eCAL Components.
	CNone       uint = 0x000
	CPublisher  uint = 0x001
	CSubscriber uint = 0x002
	CService    uint = 0x004
	CMonitoring uint = 0x008
	CLogging    uint = 0x010
	CTimeSync   uint = 0x020
	CDefault    uint = CPublisher | CSubscriber | CService | CLogging | CTimeSync
	CAll        uint = CPublisher | CSubscriber | CService | CMonitoring | CLogging | CTimeSync
)

type ProcessSeverity uint8

const (
	ProcSevUnknown  ProcessSeverity = C.process_severity_unknown
	ProcSevHealthy  ProcessSeverity = C.process_severity_healthy
	ProcSevWarning  ProcessSeverity = C.process_severity_warning
	ProcSevCritical ProcessSeverity = C.process_severity_critical
	ProcSevFailed   ProcessSeverity = C.process_severity_failed
)

type ProcessSeverityLevel uint8

const (
	ProcSevLevel1 ProcessSeverityLevel = 1
	ProcSevLevel2 ProcessSeverityLevel = 2
	ProcSevLevel3 ProcessSeverityLevel = 3
	ProcSevLevel4 ProcessSeverityLevel = 4
	ProcSevLevel5 ProcessSeverityLevel = 5
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

func GetVersionString() string {
	return C.GoString(C.GetVersionString())
}

func GetVersionDateString() string {
	return C.GoString(C.GetVersionDateString())
}

func GetVersion() C.struct_version {
	return C.GetVersion()
}

func Initialize(config Config, unitName string, components uint) bool {
	unitNameC := C.CString(unitName)

	defer C.free(unsafe.Pointer(unitNameC))

	return bool(C.Initialize(config.config, unitNameC, C.uint(components)))
}

func Finalize() bool {
	return bool(C.Finalize())
}

func IsInitialized() bool {
	return bool(C.IsInitialized())
}

func IsComponentInitialized(component uint) bool {
	return bool(C.IsComponentInitialized(C.uint(component)))
}

func Ok() bool {
	return bool(C.Ok())
}

// Set the processes state info.
func SetState(severity ProcessSeverity, level ProcessSeverityLevel, state string) {
	C.SetState(C.int(severity), C.int(level), C.CString(state))
}

// Signals an event to a local process to cause eCAL::Ok() to return false.
func ShutdownLocalEcalProcess(pid int) {
	C.ShutdownProcess(C.int(pid))
}
