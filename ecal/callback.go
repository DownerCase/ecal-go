package ecal

//#include <stdint.h>
import "C"

import (
	"runtime/cgo"
)

//export goCopyString
func goCopyString(handle C.uintptr_t, data *C.char) {
	d := cgo.Handle(handle).Value().(*string)
	*d = C.GoString(data)
}
