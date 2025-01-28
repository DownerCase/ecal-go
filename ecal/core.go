package ecal

// #cgo LDFLAGS: -lecal_core
// #include "core.h"
// #include "cgo_config.h"
// #include <stdlib.h>
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

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
