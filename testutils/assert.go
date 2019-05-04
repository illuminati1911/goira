package testutils

import "testing"

type Assert struct {
	t *testing.T
}

func NewAssert(t *testing.T) Assert {
	return Assert{t}
}

func(as *Assert) Equals(a interface{}, b interface{}) {
	if a != b {
		as.t.Errorf("Expected %s to be %s", a, b)
	}
}

func(as *Assert) NotEquals(a interface{}, b interface{}) {
	if a == b {
		as.t.Errorf("Expected %s to not be %s", a, b)
	}
}