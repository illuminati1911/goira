package hwinterface

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lm -lpigpio -pthread -lrt
// #include <stdlib.h>
// #include "ir_interface.h"
import "C"
import (
	"errors"
	"unsafe"
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

type GPIOInterface struct {
}

func NewGPIOInterface() accontrol.HWInterface {
	return &GPIOInterface{}
}

func (gpio *GPIOInterface) SetState(newState models.ACState) error {
	pin := C.uint(27)
	data := C.CString(newState.String())
	defer C.free(unsafe.Pointer(data))
	err := C.runIR(pin, data)
	if err != 0 {
		return errors.New("Failure sending IR")
	}
	return nil
}
