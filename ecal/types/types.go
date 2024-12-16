package types

import "C"

type EntityID struct {
	EntityID  string
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
