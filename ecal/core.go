package ecal

// #cgo LDFLAGS: -lecal_core
// #include "core.h"
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

type ConfigLogging struct {
	ReceiveEnabled bool
}

type Config struct {
	Logging ConfigLogging
}

type ConfigOption func(*Config)

func NewConfig(opts ...ConfigOption) Config {
	cfg := Config{}
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

func WithLoggingReceive(r bool) ConfigOption {
	return func(c *Config) {
		c.Logging.ReceiveEnabled = r
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
	cconfig := C.struct_CConfig{
		logging: C.struct_CConfigLogging{
			receive_enabled: C.bool(config.Logging.ReceiveEnabled),
		},
	}

	unitNameC := C.CString(unitName)

	defer C.free(unsafe.Pointer(unitNameC))

	return bool(C.Initialize(&cconfig, unitNameC, C.uint(components)))
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

func SetUnitName(unitName string) bool {
	unitNameC := C.CString(unitName)

	defer C.free(unsafe.Pointer(unitNameC))

	return bool(C.SetUnitName(unitNameC))
}

func Ok() bool {
	return bool(C.Ok())
}

// TODO: Reimplement with a proper config serialization as eCAL::DumpConfig()
// is planned to be removed!
func GetConfig() string {
	var cfg string

	handle := cgo.NewHandle(&cfg)
	defer handle.Delete()
	C.GetConfig(C.uintptr_t(handle))

	return cfg
}
