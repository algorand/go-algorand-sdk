package mnemonic

import (
	"fmt"
)

var errWrongKeyLen = fmt.Errorf("key length must be %d bytes", keyLenBytes)
var errWrongMnemonicLen = fmt.Errorf("mnemonic must be %d words", mnemonicLenWords)
var errWrongChecksum = fmt.Errorf("checksum failed to validate")
