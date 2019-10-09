package crypto

import (
	"errors"
)

var errInvalidSignatureReturned = errors.New("ed25519 library returned an invalid signature")
var errMsigUnknownVersion = errors.New("unknown version != 1")
var errMsigInvalidThreshold = errors.New("invalid threshold")
var errMsigBadTxnSender = errors.New("transaction sender does not match multisig parameters")
var errMsigInvalidSecretKey = errors.New("secret key has no corresponding public identity in multisig preimage")
var errMsigMergeLessThanTwo = errors.New("cannot merge fewer than two multisig transactions")
var errMsigMergeKeysMismatch = errors.New("multisig parameters do not match")
var errMsigMergeInvalidDups = errors.New("mismatched duplicate signatures")
var errLsigInvalidSignature = errors.New("invalid logicsig signature")
var errLsigInvalidProgram = errors.New("invalid logicsig program")
var errLsigEmptyMsig = errors.New("empty multisig in logicsig")
