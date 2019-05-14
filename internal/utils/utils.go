package utils


// Clamp limits value between two integers
//
// See also: https://en.cppreference.com/w/cpp/algorithm/clamp
//
func Clamp(minl, maxl, value int) int {
	return Min(Max(value, minl), maxl)
}

// Max returns larger of the two given integers. If integers are same,
// first one is returned.
//
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns smaller of the two given integers. If integers are same,
// first one is returned.
//
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Reverse reverses the bits in a byte.
//
// Example: 00100011 ==> 11000100
//
func Reverse(b byte) byte {
	b = (b & 0xF0) >> 4 | (b & 0x0F) << 4
	b = (b & 0xCC) >> 2 | (b & 0x33) << 2
	b = (b & 0xAA) >> 1 | (b & 0x55) << 1
	return b
 }
