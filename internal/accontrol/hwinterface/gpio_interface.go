package hwinterface

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"

	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

// GPIOInterface handles mapping the ACState model to binary code
// and transmits the code as IR signal via GPIO output pin.
//
type GPIOInterface struct {
	mapper accontrol.Mapper
	pin    uint
}

// iRBlasterConfig is buffer data structure that
// is sent to the kernel driver at /dev/irblaster
//
// See also: https://github.com/illuminati1911/IRBlaster
//
type iRBlasterConfig struct {
	LeadingPulseWidth  uint32
	LeadingGapWidth    uint32
	OnePulseWidth      uint32
	OneGapWidth        uint32
	ZeroPulseWidth     uint32
	ZeroGapWidth       uint32
	TrailingPulseWidth uint32
	Frequency          uint32
	DCN                uint32
	DCM                uint32
	code               [0x200]byte
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
	// TODO: Enabled when IRBlaster supports other pins than 18:
	//
	//pin := C.uint(gpio.pin)
	var b [0x200]byte
	data := gpio.mapper.MapToProtocolBinaryString(&newState)
	copy(b[:], data)
	cfg := iRBlasterConfig{
		LeadingPulseWidth:  9000,
		LeadingGapWidth:    4500,
		OnePulseWidth:      560,
		OneGapWidth:        1680,
		ZeroPulseWidth:     560,
		ZeroGapWidth:       560,
		TrailingPulseWidth: 560,
		Frequency:          38000,
		DCN:                50,
		DCM:                100,
		code:               b,
	}
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, cfg)
	if err != nil {
		return errors.New("could not generate IRBlaster transmission buffer")
	}

	f, err := os.OpenFile("/dev/irblaster", os.O_WRONLY, 0600)
	if err != nil {
		return errors.New("could not open kernel driver file descriptor")
	}
	defer f.Close()
	f.Write(buf.Bytes())
	return nil
}
