package ecal

// #cgo LDFLAGS: -lecal_core
// #include "core.h"
// #include <stdlib.h>
import "C"
import "unsafe"

const (
	// eCAL Components
	C_None       uint = 0x000
	C_Publisher  uint = 0x001
	C_Subscriber uint = 0x002
	C_Service    uint = 0x004
	C_Monitoring uint = 0x008
	C_Logging    uint = 0x010
	C_TimeSync   uint = 0x020
	C_Default    uint = C_Publisher | C_Subscriber | C_Service | C_Logging | C_TimeSync
	C_All        uint = C_Publisher | C_Subscriber | C_Service | C_Monitoring | C_Logging | C_TimeSync
)

func NewConfig() C.struct_config {
	return C.struct_config{}
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

func Initialize(config C.struct_config, unit_name string, components uint) int {
	unit_c := C.CString(unit_name)
	defer C.free(unsafe.Pointer(unit_c))
	return int(C.Initialize(&config, unit_c, C.uint(components)))
}

func Finalize() int {
	return int(C.Finalize())
}
func IsInitialized(component uint) bool {
	return bool(C.IsInitialized(C.uint(component)))
}
func SetUnitName(unit_name string) bool {
	unit_c := C.CString(unit_name)
	defer C.free(unsafe.Pointer(unit_c))
	return bool(C.SetUnitName(unit_c))
}
func Ok() bool {
	return bool(C.Ok())
}
