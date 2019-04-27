package models

import (
	"github.com/illuminati1911/goira/internal/utils"
	"fmt"
)

type ACState struct {
	Temperature *int  `json:"temp"`
	WindLevel   *int  `json:"wind"` // Ignored for now
	Mode        *int  `json:"mode"`
	Active      *bool `json:"active"`
}

func (ac *ACState) String() string {
	bytes := ac.stateAsBytes()
	bytes = append(bytes, ac.checksum(bytes))
	byteStr := ""
	for _, b := range bytes {
		byteStr += fmt.Sprintf("%08b", utils.Reverse(b))
	}
	fmt.Println(byteStr)
	return byteStr
}

func (ac *ACState) stateAsBytes() []byte {
	cTemp := utils.Clamp(16, 30, *ac.Temperature)
	bTemp := byte(92 + cTemp)
	cMode := utils.Clamp(1, 3, *ac.Mode)
	bMode := byte(cMode * 16)
	bActive := byte(0xC0)
	if *ac.Active {
		bActive = byte(0x00)
	}
	return []byte{
		0x56,		// Header
		bTemp,		// Temperature
		0x00,
		0x00,
		bMode,		// AC mode
		bActive,	// On or off
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}
}

func (ac *ACState) checksum(bytes []byte) byte {
	var checksum byte = 0x00
	for _, b := range bytes {
		checksum += (b & 0xF0) >> 4
		checksum += b & 0x0F
	}
	return checksum
}