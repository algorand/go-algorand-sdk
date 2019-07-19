package types

import (
	"errors"
	"golang.org/x/crypto/ed25519"
)

// Signature is an ed25519 signature
type Signature [ed25519.SignatureSize]byte

// ToBytes converts a signature to a byte representation
func (s Signature) ToBytes() []byte {
	return s[:]
}

// MultisigSubsig contains a single public key and, optionally, a signature
type MultisigSubsig struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Key ed25519.PublicKey `codec:"pk"`
	Sig Signature         `codec:"s"`
}

// MultisigSig holds multiple Subsigs, as well as threshold and version info
type MultisigSig struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Version   uint8            `codec:"v"`
	Threshold uint8            `codec:"thr"`
	Subsigs   []MultisigSubsig `codec:"subsig"`
}

// MakeSignature converts data into a Signature and checks the size.
func MakeSignature(data []byte) (s Signature, err error) {
	n := copy(s[:], data)
	if n != len(s) {
		err = errors.New("ed25519 library returned an invalid signature")
		return
	}
	return
}
