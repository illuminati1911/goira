package testutils

import "testing"

// Assert is a type containing multiple different assertation tools
// for unit testing.
//
type Assert struct {
	t *testing.T
}

// NewAssert crates new instance of Assert with testing.T pointer.
//
func NewAssert(t *testing.T) Assert {
	return Assert{t}
}

// Equals determines whether two given parameters are equal
//
func(as *Assert) Equals(a interface{}, b interface{}) {
	if a != b {
		as.t.Errorf("Expected %s to be %s", a, b)
	}
}

// NotEquals determines whether two given parameters are not equal
//
func(as *Assert) NotEquals(a interface{}, b interface{}) {
	if a == b {
		as.t.Errorf("Expected %s to not be %s", a, b)
	}
}