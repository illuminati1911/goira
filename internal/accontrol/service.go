package accontrol

type Service interface {
	SetTemperature(temp int) error
	SetWindLevel(level int) error
	TurnOn() error
	TurnOff() error
}
