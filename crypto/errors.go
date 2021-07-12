package crypto

import (
	"errors"
)

var errInvalidSignatureReturned = errors.New("ed25519 library returned an invalid signature")
var errInvalidPrivateKey = errors.New("invalid private key")
var errMsigUnknownVersion = errors.New("unknown version != 1")
var errMsigInvalidThreshold = errors.New("invalid threshold")
var errMsigInvalidSecretKey = errors.New("secret key has no corresponding public identity in multisig preimage")
var errMsigMergeLessThanTwo = errors.New("cannot merge fewer than two multisig transactions")
var errMsigMergeKeysMismatch = errors.New("multisig parameters do not match")
var errMsigMergeInvalidDups = errors.New("mismatched duplicate signatures")
var errLsigTooManySignatures = errors.New("logicsig has too many signatures, at most one of Sig or Msig may be defined")
var errLsigNoSignature = errors.New("logicsig is not delegated")
var errLsigNoSigningKey = errors.New("signed logicsig does have a valid signing key, try signing again")
var errLsigInvalidSignature = errors.New("invalid logicsig signature")
var errLsigInvalidProgram = errors.New("invalid logicsig program")
var errLsigEmptyMsig = errors.New("empty multisig in logicsig")
