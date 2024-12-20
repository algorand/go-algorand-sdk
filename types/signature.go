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
//
// LogicSig cannot sign transactions in all cases.  Instead, use LogicSigAccount as a safe, general purpose signing mechanism.  Since LogicSig does not track the provided signature's public key, LogicSig cannot sign transactions when delegated to a non-multisig account _and_ the sender is not the delegating account.
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

/* Classical signatures */
type ed25519Signature [64]byte
type ed25519PublicKey [32]byte

// A HeartbeatProof is functionally equivalent to a OneTimeSignature (see below), but it has
// been cleaned up for use as a transaction field in heartbeat transactions.
//
// A OneTimeSignature is a cryptographic signature that is produced a limited
// number of times and provides forward integrity.
//
// Specifically, a OneTimeSignature is generated from an ephemeral secret. After
// some number of messages is signed under a given OneTimeSignatureIdentifier
// identifier, the corresponding secret is deleted. This prevents the
// secret-holder from signing a contradictory message in the future in the event
// of a secret-key compromise.
type HeartbeatProof struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Sig is a signature of msg under the key PK.
	Sig ed25519Signature `codec:"s"`
	PK  ed25519PublicKey `codec:"p"`

	// PK2 is used to verify a two-level ephemeral signature.
	PK2 ed25519PublicKey `codec:"p2"`
	// PK1Sig is a signature of OneTimeSignatureSubkeyOffsetID(PK, Batch, Offset) under the key PK2.
	PK1Sig ed25519Signature `codec:"p1s"`
	// PK2Sig is a signature of OneTimeSignatureSubkeyBatchID(PK2, Batch) under the master key (OneTimeSignatureVerifier).
	PK2Sig ed25519Signature `codec:"p2s"`
}
