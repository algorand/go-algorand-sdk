package types

import (
	"fmt"
)

var errWrongAddressLen = fmt.Errorf("decoded address is the wrong length, should be %d bytes", hashLenBytes+checksumLenBytes)
var errWrongChecksum = fmt.Errorf("address checksum is incorrect, did you copy the address correctly?")
