package crypto

import "github.com/algorand/go-algorand-sdk/v2/types"

// Ed25519PublicKeySize is the size in bytes of an ed25519 public key
const Ed25519PublicKeySize = 32

// Ed25519PublicKey represents a 32 byte ed25519 public key.
type Ed25519PublicKey [Ed25519PublicKeySize]byte

// Falcon1024PublicKeySize is the size in bytes of a falcon1024 public key
const Falcon1024PublicKeySize = 1793

// Falcon1024PublicKey represents a 1793 byte falcon1024 public key.
type Falcon1024PublicKey [Falcon1024PublicKeySize]byte

// Ed25519Signer represents the ability to perform ed25519 signatures on behalf
// of some public key
type Ed25519Signer interface {
	// Ed25519Sign signs the given bytes with an ed25519 signature
	Ed25519Sign(toBeSigned []byte) ([]byte, error)
	// Ed25519PublicKey returns the public key that should be used to verify
	// the signatures performed by this signer
	Ed25519PublicKey() Ed25519PublicKey
}

// Falcon1024Signer represents the ability to perform falcon1024 signatures on
// behalf of some public key
//
// Signers for non-canonical accounts should also implement Falcon1024Salted
type Falcon1024Signer interface {
	// Falcon1024Sign signs the given bytes with a falcon1024 signature
	Falcon1024Sign(toBeSigned []byte) ([]byte, error)
	// Falcon1024PublicKey returns the public key that should be used to verify
	// the signatures performed by this signer
	Falcon1024PublicKey() Falcon1024PublicKey
}

// Falcon1024Salted equips a signer with the ability to specify a custom (maybe
// non-canonical) salt
type Falcon1024Salted interface {
	// Falcon1024Salt returns the (maybe non-canonical) salt that identifies the
	// account selected for this signer
	Falcon1024Salt() types.PQAddressSalt
}
