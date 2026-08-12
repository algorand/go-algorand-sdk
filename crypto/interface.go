package crypto

// Ed25519PublicKey represents a 32 byte ed25519 public key.
type Ed25519PublicKey [32]byte

// Ed25519Signer represents the ability to perform ed25519 signatures on behalf
// of some public key
type Ed25519Signer interface {
	// Ed25519Sign signs the given bytes with an ed25519 signature
	Ed25519Sign(toBeSigned []byte) ([]byte, error)
	// Ed25519PublicKey returns the public key that should be used to verify
	// the signatures performed by this signer
	Ed25519PublicKey() Ed25519PublicKey
}
