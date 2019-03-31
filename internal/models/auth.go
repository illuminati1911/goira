package models

import (
	"time"
)

type Token struct {
	Name    string
	Value   string
	Expires time.Time
}

type Password = string
