package ecal

// #cgo LDFLAGS: -lecal_core
// #include "ecal_go_core.h"
import "C"

const (
	None       C.uint = 0x000
	Publisher         = 0x001
	Subscriber        = 0x002
	Service           = 0x004
	Monitoring        = 0x008
	Logging           = 0x010
	TimeSync          = 0x020
	Default           = Publisher | Subscriber | Service | Logging | TimeSync
	All               = Publisher | Subscriber | Service | Monitoring | Logging | TimeSync
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
