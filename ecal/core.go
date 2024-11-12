package ecal

// #cgo LDFLAGS: -lecal_core
// #include "core.h"
import "C"

const (
	// eCAL Components
	C_None       C.uint = 0x000
	C_Publisher         = 0x001
	C_Subscriber        = 0x002
	C_Service           = 0x004
	C_Monitoring        = 0x008
	C_Logging           = 0x010
	C_TimeSync          = 0x020
	C_Default           = C_Publisher | C_Subscriber | C_Service | C_Logging | C_TimeSync
	C_All               = C_Publisher | C_Subscriber | C_Service | C_Monitoring | C_Logging | C_TimeSync
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

func Initialize(config C.struct_config, unit_name string, components C.uint) int {
	return int(C.Initialize(&config, C.CString(unit_name), components))
}

func Finalize() int {
	return int(C.Finalize())
}
func IsInitialized(component C.uint) bool {
	return bool(C.IsInitialized(component))
}
func SetUnitName(unit_name string) bool {
	return bool(C.SetUnitName(C.CString(unit_name)))
}
func Ok() bool {
	return bool(C.Ok())
}
