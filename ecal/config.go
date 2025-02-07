package ecal

// #cgo LDFLAGS: -lecal_core
// #include "cgo_config.h"
// #include <stdlib.h>
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type Config struct {
	config unsafe.Pointer
}

type ConfigOption func(*Config)

func NewConfig(opts ...ConfigOption) Config {
	cfg := Config{
		config: C.NewConfig(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

func (c *Config) Delete() {
	C.DeleteConfig(c.config)

	c.config = nil
}

// Enable/Disable printing eCAL logs to stderr.
func WithConsoleLogging(enable bool) ConfigOption {
	return func(c *Config) {
		C.ConfigLoggingConsole(c.config, C.bool(enable))
	}
}

// Enable all log levels when printing eCAL logs to stderr.
func WithConsoleLogAll() ConfigOption {
	return func(c *Config) {
		C.ConfigLoggingConsoleAll(c.config)
	}
}

// Enable/Disable sending eCAL logs over UDP.
func WithUDPLogging(enable bool) ConfigOption {
	return func(c *Config) {
		C.ConfigLoggingUdp(c.config, C.bool(enable))
	}
}

// Enable all log levels when sending eCAL logs over UDP.
func WithUDPLogAll() ConfigOption {
	return func(c *Config) {
		C.ConfigLoggingUdpAll(c.config)
	}
}

func WithUDPLogReceiving(enable bool) ConfigOption {
	return func(c *Config) {
		C.ConfigLoggingUdpReceive(c.config, C.bool(enable))
	}
}

// WARNING: These functions will return the default values before Initialize has been called

func GetLoadedConfigFilePath() string {
	var cfg string

	handle := cgo.NewHandle(&cfg)
	defer handle.Delete()
	C.ConfigGetLoadedFilePath(C.uintptr_t(handle))

	return cfg
}

func PublisherShmEnabled() bool {
	return bool(C.ConfigPubShmEnabled())
}

func PublisherUDPEnabled() bool {
	return bool(C.ConfigPubUdpEnabled())
}

func PublisherTCPEnabled() bool {
	return bool(C.ConfigPubTcpEnabled())
}

func SubscriberShmEnabled() bool {
	return bool(C.ConfigSubShmEnabled())
}

func SubscriberUDPEnabled() bool {
	return bool(C.ConfigSubUdpEnabled())
}

func SubscriberTCPEnabled() bool {
	return bool(C.ConfigSubTcpEnabled())
}
