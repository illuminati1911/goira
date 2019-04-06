package models

type ACState struct {
	Temperature *int  `json:"temp"`
	WindLevel   *int  `json:"wind"`
	Active      *bool `json:"active"`
}
