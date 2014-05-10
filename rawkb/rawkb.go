package rawkb

/*
#include <stdlib.h>
#include <unistd.h>
#include "test.h"
#cgo CFLAGS:-Wno-error
*/
import "C"
import "unsafe"

func SetupKeyboard() int {
	i := C.setupKeyboard()
	return int(i)
}

func RestoreKeyboard() {
	C.restoreKeyboard()
}

func ReadOnce() (uint16, uint8) {
	var a uint16
	i := C.read(0, unsafe.Pointer(&a), 1)
	return a, uint8(i)
}
