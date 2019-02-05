package crypto

import (
	"fmt"
)

var errInvalidSignatureReturned = fmt.Errorf("ed25519 library returned an invalid signature")
