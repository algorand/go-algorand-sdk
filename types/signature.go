package types

import (
	"golang.org/x/crypto/ed25519"
)

// Signature is an ed25519 signature
type Signature [ed25519.SignatureSize]byte

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

// Blank returns true iff the msig is empty. We need this instead of just
// comparing with == MultisigSig{}, because Subsigs is a slice.
func (msig MultisigSig) Blank() bool {
	if msig.Version != 0 {
		return false
	}
	if msig.Threshold != 0 {
		return false
	}
	if msig.Subsigs != nil {
		return false
	}
	return true
}

// LogicSig contains logic for validating a transaction.
// LogicSig is signed by an account, allowing delegation of operations.
// OR
// LogicSig defines a contract account.
type LogicSig struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Logic signed by Sig or Msig
	// OR hashed to be the Address of an account.
	Logic []byte `codec:"l"`

	// The signature of the account that has delegated to this LogicSig, if any
	Sig Signature `codec:"sig"`

	// The signature of the multisig account that has delegated to this LogicSig, if any
	Msig MultisigSig `codec:"msig"`

	// Args are not signed, but checked by Logic
	Args [][]byte `codec:"arg"`
}

// Blank returns true iff the lsig is empty. We need this instead of just
// comparing with == LogicSig{}, because it contains slices.
func (lsig LogicSig) Blank() bool {
	if lsig.Args != nil {
		return false
	}
	if len(lsig.Logic) != 0 {
		return false
	}
	if !lsig.Msig.Blank() {
		return false
	}
	if lsig.Sig != (Signature{}) {
		return false
	}
	return true
}
