package crypto

import (
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/types"
)

// Account holds both the public and private information associated with an
// Algorand address
type Account struct {
	PublicKey  []byte
	PrivateKey []byte
	Address    []byte
}

func init() {
	addrLen := len(types.Address{})
	pkLen := ed25519.PublicKeySize
	if addrLen != pkLen {
		panic("address and public key are different sizes")
	}
}

// GenerateAccount generates a random Account
func GenerateAccount(pk, sk []byte) string {
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

	return a.String()
}
