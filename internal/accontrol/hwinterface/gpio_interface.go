package hwinterface

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lm -lpigpio -pthread -lrt
// #include <stdlib.h>
// #include "ir_interface.h"
import "C"
import (
	"errors"
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
	"unsafe"
)

// GPIOInterface handles mapping the ACState model to binary code
// and transmits the code as IR signal via GPIO output pin.
//
type GPIOInterface struct {
	mapper accontrol.Mapper
	pin    uint
}

// NewGPIOInterface creates new instance of GPIOInterface
//
func NewGPIOInterface(mapper accontrol.Mapper, pin uint) accontrol.HWInterface {
	return &GPIOInterface{mapper, pin}
}

// SetState will convert the model to binary string and pass it
// to the irslinger.
//
func (gpio *GPIOInterface) SetState(newState models.ACState) error {
	pin := C.uint(gpio.pin)
	data := C.CString(gpio.mapper.MapToProtocolBinaryString(&newState))
	defer C.free(unsafe.Pointer(data))
	err := C.runIR(pin, data)
	if err != 0 {
		return errors.New("Failure sending IR")
	}
	return nil
}
