package crypto

import (
	"fmt"
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

func GenerateSK() []byte {
	_, sk, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	return sk
}

func GenerateAddressFromSK(sk []byte) (string, error) {
	edsk := ed25519.PrivateKey(sk)

	var a types.Address
	pk := edsk.Public()
	n := copy(a[:], []byte(pk.(ed25519.PublicKey)))
	if n != ed25519.PublicKeySize {
		return "", fmt.Errorf("generated public key has the wrong size, expected %d, got %d", ed25519.PublicKeySize, n)
	}
	return a.String(), nil
}
