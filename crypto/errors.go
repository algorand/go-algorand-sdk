package crypto

import (
	"errors"
)

var errMsigUnknownVersion = errors.New("unknown version != 1")
var errMsigInvalidThreshold = errors.New("invalid threshold")
var errMsigBadTxnSender = errors.New("transaction sender does not match multisig parameters")
var errMsigInvalidSecretKey = errors.New("secret key has no corresponding public identity in multisig preimage")
var errMsigMergeLessThanTwo = errors.New("cannot merge fewer than two multisig transactions")
var errMsigMergeKeysMismatch = errors.New("multisig parameters do not match")
var errMsigMergeInvalidDups = errors.New("mismatched duplicate signatures")
