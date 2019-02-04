package crypto

import (
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/types"
)

// KeyPair holds both the public and private information associated with an
// Algorand address
type KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	Address    types.Address
}

func init() {
	addrLen := len(types.Address{})
	pkLen := ed25519.PublicKeySize
	if addrLen != pkLen {
		panic("address and public key are different sizes")
	}
}

// GenerateKeyPair generates a random KeyPair
func GenerateKeyPair() (kp KeyPair) {
	// Generate an ed25519 keypair. This should never fail
	pk, sk, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	// Convert the public key to an address
	var a types.Address
	n := copy(a[:], pk)
	if n != ed25519.PublicKeySize {
		panic("generated public key is the wrong size")
	}

	// Build the returned keypair
	kp.PublicKey = pk
	kp.PrivateKey = sk
	kp.Address = a
	return
}
