package crypto

import "github.com/algorand/go-algorand-sdk/v2/types"

// Ed25519PublicKeySize is the size in bytes of an ed25519 public key
const Ed25519PublicKeySize = 32

// Ed25519PublicKey represents a 32 byte ed25519 public key.
type Ed25519PublicKey [Ed25519PublicKeySize]byte

// Ed25519Signer represents the ability to perform ed25519 signatures on behalf
// of some public key
type Ed25519Signer interface {
	// Ed25519Sign signs the given bytes with an ed25519 signature
	Ed25519Sign(toBeSigned []byte) ([]byte, error)
	// Ed25519PublicKey returns the public key that should be used to verify
	// the signatures performed by this signer
	Ed25519PublicKey() Ed25519PublicKey
}

// PQSigner represents the ability to perform pq signatures on
// behalf of some public key
//
// Signers for non-canonical accounts should also implement Falcon1024Salted
type PQSigner interface {
	// PQSign signs the given bytes with a pq signature
	PQSign(toBeSigned []byte) ([]byte, error)
	// PQPublicKey returns the public key that should be used to verify
	// the signatures performed by this signer
	PQPublicKey() []byte
	// PQScheme returns the identifier for the post-quantum scheme used by this
	// signer
	PQScheme() types.PQScheme
}

// PQSalted equips a signer with the ability to specify a custom (maybe
// non-canonical) salt
type PQSalted interface {
	// PQSalt returns the (maybe non-canonical) salt that identifies the
	// account selected for this signer
	PQSalt() types.PQAddressSalt
}
