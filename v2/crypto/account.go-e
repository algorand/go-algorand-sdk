package crypto

import (
	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/types"
)

// prefix for multisig transaction signing
const msigAddrPrefix = "MultisigAddr"

// Account holds both the public and private information associated with an
// Algorand address
type Account struct {
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

// GenerateAccount generates a random Account
func GenerateAccount() (kp Account) {
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

	// Build the account
	kp.PublicKey = pk
	kp.PrivateKey = sk
	kp.Address = a
	return
}

/* Multisig Support */

// MultisigAccount is a convenience type for holding multisig preimage data
type MultisigAccount struct {
	// Version is the version of this multisig
	Version uint8
	// Threshold is how many signatures are needed to fully sign as this address
	Threshold uint8
	// Pks is an ordered list of public keys that could potentially sign a message
	Pks []ed25519.PublicKey
}

// MultisigAccountWithParams creates a MultisigAccount with the given parameters
func MultisigAccountWithParams(version uint8, threshold uint8, addrs []types.Address) (ma MultisigAccount, err error) {
	ma.Version = version
	ma.Threshold = threshold
	ma.Pks = make([]ed25519.PublicKey, len(addrs))
	for i := 0; i < len(addrs); i++ {
		ma.Pks[i] = addrs[i][:]
	}
	err = ma.Validate()
	return
}

// MultisigAccountFromSig is a convenience method that creates an account
// from a sig in a signed tx. Useful for getting addresses from signed msig txs, etc.
func MultisigAccountFromSig(sig types.MultisigSig) (ma MultisigAccount, err error) {
	ma.Version = sig.Version
	ma.Threshold = sig.Threshold
	ma.Pks = make([]ed25519.PublicKey, len(sig.Subsigs))
	for i := 0; i < len(sig.Subsigs); i++ {
		c := make([]byte, len(sig.Subsigs[i].Key))
		copy(c, sig.Subsigs[i].Key)
		ma.Pks[i] = c
	}
	err = ma.Validate()
	return
}

// Address takes this multisig preimage data, and generates the corresponding identifying
// address, committing to the exact group, version, and public keys that it requires to sign.
// Hash("MultisigAddr" || version uint8 || threshold uint8 || PK1 || PK2 || ...)
func (ma MultisigAccount) Address() (addr types.Address, err error) {
	// See go-algorand/crypto/multisig.go
	err = ma.Validate()
	if err != nil {
		return
	}
	buffer := append([]byte(msigAddrPrefix), byte(ma.Version), byte(ma.Threshold))
	for _, pki := range ma.Pks {
		buffer = append(buffer, pki[:]...)
	}
	return sha512.Sum512_256(buffer), nil
}

// Validate ensures that this multisig setup is a valid multisig account
func (ma MultisigAccount) Validate() (err error) {
	if ma.Version != 1 {
		err = errMsigUnknownVersion
		return
	}
	if ma.Threshold == 0 || len(ma.Pks) == 0 || int(ma.Threshold) > len(ma.Pks) {
		err = errMsigInvalidThreshold
		return
	}
	return
}

// Blank return true if MultisigAccount is empty
// struct containing []ed25519.PublicKey cannot be compared
func (ma MultisigAccount) Blank() bool {
	if ma.Version != 0 {
		return false
	}
	if ma.Threshold != 0 {
		return false
	}
	if ma.Pks != nil {
		return false
	}
	return true
}

// LogicSigAddress returns contract (escrow) address
func LogicSigAddress(lsig types.LogicSig) types.Address {
	toBeSigned := programToSign(lsig.Logic)
	checksum := sha512.Sum512_256(toBeSigned)

	var addr types.Address
	n := copy(addr[:], checksum[:])
	if n != ed25519.PublicKeySize {
		panic(fmt.Sprintf("Generated public key has length of %d, expected %d", n, ed25519.PublicKeySize))
	}
	return addr
}
