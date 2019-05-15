package models

type ACState struct {
	Temperature *int  `json:"temp"`
	WindLevel   *int  `json:"wind"` // Ignored for now
	Mode        *int  `json:"mode"`
	Active      *bool `json:"active"`
}
