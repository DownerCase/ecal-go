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

type LogLevel uint8

type LogMessage struct {
	Time      int64
	Host      string
	Process   string
	Unit_name string
	Content   string
	Pid       int32
	Level     LogLevel
}
