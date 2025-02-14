package ecallog

// #cgo CPPFLAGS: -I${SRCDIR}/..
// #include "logging.h"
// void GoLog(enum CLogLevel level, _GoString_ msg) {
//   Log(level, _GoStringPtr(msg), _GoStringLen(msg));
// }
//// C preamble.
import "C"

import (
	"fmt"
	"runtime/cgo"
	"strconv"
)

type Level uint8

type LogMessage struct {
	Time     int64
	Host     string
	Process  string
	UnitName string
	Content  string
	Pid      int32
	Level    Level
}

type LogMessages struct {
	Messages []LogMessage
}

const (
	LogLevelNone   Level = C.log_level_none
	LogLevelAll    Level = C.log_level_all
	LogLevelInfo   Level = C.log_level_info
	LogLevelWarn   Level = C.log_level_warning
	LogLevelError  Level = C.log_level_error
	LogLevelFatal  Level = C.log_level_fatal
	LogLevelDebug  Level = C.log_level_debug1
	LogLevelDebug1 Level = C.log_level_debug1
	LogLevelDebug2 Level = C.log_level_debug2
	LogLevelDebug3 Level = C.log_level_debug3
	LogLevelDebug4 Level = C.log_level_debug4
)

func (l Level) String() string {
	switch l {
	case LogLevelNone:
		return "None"
	case LogLevelAll:
		return "All"
	case LogLevelInfo:
		return "Info"
	case LogLevelWarn:
		return "Warn"
	case LogLevelError:
		return "Error"
	case LogLevelFatal:
		return "Fatal"
	case LogLevelDebug1:
		return "Debug1"
	case LogLevelDebug2:
		return "Debug2"
	case LogLevelDebug3:
		return "Debug3"
	case LogLevelDebug4:
		return "Debug4"
	default:
		return strconv.FormatUint(uint64(l), 10)
	}
}

func Log(level Level, a ...any) {
	C.GoLog(uint32(level), fmt.Sprint(a...))
}

func Logf(level Level, format string, a ...any) {
	C.GoLog(uint32(level), fmt.Sprintf(format, a...))
}

func Error(a ...any) {
	Log(LogLevelError, a...)
}

func Errorf(format string, a ...any) {
	Logf(LogLevelError, format, a...)
}

func Warn(a ...any) {
	Log(LogLevelWarn, a...)
}

func Warnf(format string, a ...any) {
	Logf(LogLevelWarn, format, a...)
}

func Info(a ...any) {
	Log(LogLevelInfo, a...)
}

func Infof(format string, a ...any) {
	Logf(LogLevelInfo, format, a...)
}

func Debug(a ...any) {
	Log(LogLevelDebug, a...)
}

func Debugf(format string, a ...any) {
	Logf(LogLevelDebug, format, a...)
}

func GetLogging() LogMessages {
	var logs LogMessages
	handle := cgo.NewHandle(&logs)
	C.GetLogging(C.uintptr_t(handle))
	handle.Delete()

	return logs
}
