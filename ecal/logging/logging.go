package logging

// #cgo CPPFLAGS: -I${SRCDIR}/..
// #include <ecal/ecal_log_level.h>
// #include "logging.h"
// void GoLog(enum eCAL_Logging_eLogLevel level, _GoString_ msg) {
//   Log(level, _GoStringPtr(msg), _GoStringLen(msg));
// }
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

const (
	LevelNone   Level = C.log_level_none
	LevelAll    Level = C.log_level_all
	LevelInfo   Level = C.log_level_info
	LevelWarn   Level = C.log_level_warning
	LevelError  Level = C.log_level_error
	LevelFatal  Level = C.log_level_fatal
	LevelDebug  Level = C.log_level_debug1
	LevelDebug1 Level = C.log_level_debug1
	LevelDebug2 Level = C.log_level_debug2
	LevelDebug3 Level = C.log_level_debug3
	LevelDebug4 Level = C.log_level_debug4
)

func (l Level) String() string {
	switch l {
	case LevelNone:
		return "None"
	case LevelAll:
		return "All"
	case LevelInfo:
		return "Info"
	case LevelWarn:
		return "Warn"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "Fatal"
	case LevelDebug1:
		return "Debug1"
	case LevelDebug2:
		return "Debug2"
	case LevelDebug3:
		return "Debug3"
	case LevelDebug4:
		return "Debug4"
	default:
		return strconv.FormatUint(uint64(l), 10)
	}
}

type Logging struct {
	Messages []LogMessage
}

func GetLogging() Logging {
	var logs Logging
	handle := cgo.NewHandle(&logs)
	C.GetLogging(C.uintptr_t(handle))
	handle.Delete()
	return logs
}

func Log(level Level, a ...any) {
	C.GoLog(uint32(level), fmt.Sprint(a...))
}

func Logf(level Level, format string, a ...any) {
	C.GoLog(uint32(level), fmt.Sprintf(format, a...))
}

func Error(a ...any) {
	Log(LevelError, a...)
}

func Errorf(format string, a ...any) {
	Logf(LevelError, format, a...)
}

func Warn(a ...any) {
	Log(LevelWarn, a...)
}

func Warnf(format string, a ...any) {
	Logf(LevelWarn, format, a...)
}

func Info(a ...any) {
	Log(LevelInfo, a...)
}

func Infof(format string, a ...any) {
	Logf(LevelInfo, format, a...)
}

func Debug(a ...any) {
	Log(LevelDebug, a...)
}

func Debugf(format string, a ...any) {
	Logf(LevelDebug, format, a...)
}
