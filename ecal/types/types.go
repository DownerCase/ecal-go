package types

import "C"

type EntityId struct {
	Entity_id  string
	Process_id int32
	Host_name  string
}

type TopicId struct {
	Topic_id   EntityId
	Topic_name string
}

type DataType struct {
	Name       string
	Encoding   string
	Descriptor []byte
}
