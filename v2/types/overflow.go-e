package types

// Overflow provides a set of helper function to prevent overflows

// OAdd16 adds 2 uint16 values with overflow detection
func OAdd16(a uint16, b uint16) (res uint16, overflowed bool) {
	res = a + b
	overflowed = res < a
	return
}

// OAdd adds 2 values with overflow detection
func OAdd(a uint64, b uint64) (res uint64, overflowed bool) {
	res = a + b
	overflowed = res < a
	return
}

// OSub subtracts b from a with overflow detection
func OSub(a uint64, b uint64) (res uint64, overflowed bool) {
	res = a - b
	overflowed = res > a
	return
}

// OMul multiplies 2 values with overflow detection
func OMul(a uint64, b uint64) (res uint64, overflowed bool) {
	if b == 0 {
		return 0, false
	}

	c := a * b
	if c/b != a {
		return 0, true
	}
	return c, false
}
