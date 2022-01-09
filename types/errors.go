package types

import (
	"fmt"
)

var errWrongAddressByteLen = fmt.Errorf("encoding address is the wrong length, should be %d bytes", hashLenBytes)
var errWrongAddressLen = fmt.Errorf("decoded address is the wrong length, should be %d bytes", hashLenBytes+checksumLenBytes)
var errWrongChecksum = fmt.Errorf("address checksum is incorrect, did you copy the address correctly?")
