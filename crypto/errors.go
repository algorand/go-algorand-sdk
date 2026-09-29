package crypto

import (
	"errors"

	"github.com/algorand/go-algorand-sdk/v2/internal/signing"
)

// ErrInvalidSignatureReturned is returned when an Ed25519Signer produces a
// signature whose length is not ed25519.SignatureSize.
var ErrInvalidSignatureReturned = signing.ErrInvalidSignatureReturned

var errInvalidPrivateKey = errors.New("invalid private key")
var errLsigTooManySignatures = errors.New("logicsig has too many signatures, at most one of Sig, Msig, LMsig or PQsig may be defined")
var errLsigInvalidSignature = errors.New("invalid logicsig signature")
var errLsigNoPublicKey = errors.New("missing public key of delegated logicsig")
var errLsigInvalidPublicKey = errors.New("public key does not match logicsig signature")
var errLsigAccountPublicKeyNotNeeded = errors.New("a public key for the signer was provided when none was expected")
