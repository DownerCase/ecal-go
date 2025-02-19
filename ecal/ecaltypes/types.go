package ecaltypes

import (
	"unsafe"
)

type EntityID struct {
	EntityID  uint64
	ProcessID int32
	HostName  string
}

type TopicID struct {
	TopicID   EntityID
	TopicName string
}

type DataType struct {
	Name       string
	Encoding   string
	Descriptor []byte
}

type Subscriber struct {
	Callback func(unsafe.Pointer, int)
}
