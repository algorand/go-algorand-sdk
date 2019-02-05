package types

import (
	"golang.org/x/crypto/ed25519"
)

// Signature is an ed25519 signature
type Signature [ed25519.SignatureSize]byte

// MultisigSubsig contains a single public key and, optionally, a signature
type MultisigSubsig struct {
	Key ed25519.PublicKey `codec:"pk"`
	Sig Signature         `codec:"s"`
}

// MultisigSig holds multiple Subsigs, as well as threshold and version info
type MultisigSig struct {
	Version   uint8            `codec:"v"`
	Threshold uint8            `codec:"thr"`
	Subsigs   []MultisigSubsig `codec:"subsig"`
}
