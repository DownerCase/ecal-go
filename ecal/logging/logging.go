package logging

// #include <ecal/ecal_log_level.h>
// #include "logging.h"
// void GoLog(enum eCAL_Logging_eLogLevel level, _GoString_ msg) {
//   Log(level, _GoStringPtr(msg), _GoStringLen(msg));
// }
import "C"
import "fmt"

type Level uint8

const (
	LevelNone   Level = C.log_level_none
	LevelAll          = C.log_level_all
	LevelInfo         = C.log_level_info
	LevelWarn         = C.log_level_warning
	LevelError        = C.log_level_error
	LevelFatal        = C.log_level_fatal
	LevelDebug        = C.log_level_debug1
	LevelDebug1       = C.log_level_debug1
	LevelDebug2       = C.log_level_debug2
	LevelDebug3       = C.log_level_debug3
	LevelDebug4       = C.log_level_debug4
)

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

func SetConsoleFilter(levels Level) {
	C.SetConsoleFilter(C.eCAL_Logging_Filter(levels))
}

func SetUdpFilter(levels Level) {
	C.SetUDPFilter(C.eCAL_Logging_Filter(levels))
}

func SetFileFilter(levels Level) {
	C.SetFileFilter(C.eCAL_Logging_Filter(levels))
}
