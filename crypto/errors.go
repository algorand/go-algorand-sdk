package crypto

import (
	"errors"
)

// ErrInvalidSignatureReturned is returned when an Ed25519Signer produces a
// signature whose length is not ed25519.SignatureSize.
var ErrInvalidSignatureReturned = errors.New("ed25519 library returned an invalid signature")

// ErrNilPQSigner is returned when a post-quantum signing operation is given no
// signer to sign with.
var ErrNilPQSigner = errors.New("pq signer cannot be nil")

var errInvalidPrivateKey = errors.New("invalid private key")
var errMsigUnknownVersion = errors.New("unknown version != 1")
var errMsigInvalidThreshold = errors.New("invalid threshold")
var errMsigInvalidSecretKey = errors.New("secret key has no corresponding public identity in multisig preimage")
var errMsigMergeLessThanTwo = errors.New("cannot merge fewer than two multisig transactions")
var errMsigMergeKeysMismatch = errors.New("multisig parameters do not match")
var errMsigMergeInvalidDups = errors.New("mismatched duplicate signatures")
var errMsigMergeAuthAddrMismatch = errors.New("mismatched AuthAddrs")
var errLsigTooManySignatures = errors.New("logicsig has too many signatures, at most one of Sig, Msig, LMsig or PQsig may be defined")
var errLsigInvalidSignature = errors.New("invalid logicsig signature")
var errLsigNoPublicKey = errors.New("missing public key of delegated logicsig")
var errLsigInvalidPublicKey = errors.New("public key does not match logicsig signature")
var errLsigEmptyMsig = errors.New("empty multisig in logicsig")
var errLsigAccountPublicKeyNotNeeded = errors.New("a public key for the signer was provided when none was expected")
