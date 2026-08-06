package crypto

import (
	"crypto/sha512"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// pqAddressPrefix is prepended when deriving a post-quantum account address.
var pqAddressPrefix = []byte("PQA")

// pqProgramPrefix is prepended to a logic program when computing the bytes a
// post-quantum scheme signs for a delegated LogicSig.
var pqProgramPrefix = []byte("PQProgram")

const Falcon1024PublicKeySize = 1793
const Falcon1024PrivateKeySize = 2305

// Falcon1024Account holds both the public and private information associated with a
// falcon address.
type Falcon1024Account struct {
	PublicKey  [Falcon1024PublicKeySize]byte
	PrivateKey [Falcon1024PrivateKeySize]byte
	Salt       types.PQAddressSalt
}

// Address returns the account address for the given Falcon1024Account.
// Hash("PQA"  || scheme ||  salt || publicKey)
func (pqa Falcon1024Account) Address() (addr types.Address) {
	buf := make([]byte, 0, len(pqAddressPrefix)+len(types.PQSchemeFalcon1024)+1+len(pqa.PublicKey))
	buf = append(buf, pqAddressPrefix...)
	buf = append(buf, types.PQSchemeFalcon1024[:]...)
	buf = append(buf, uint8(pqa.Salt))
	buf = append(buf, pqa.PublicKey[:]...)

	digest := sha512.Sum512_256(buf)

	copy(addr[:], digest[:])

	return
}

// Validate returns an error if the given Falcon1024Account address could be interpreted as an ed25519 public key
func (pqa Falcon1024Account) Validate() error {
	addr := pqa.Address()

	if IsEdwards25519Point([]byte(addr[:])) {
		return fmt.Errorf("Account address overlaps with a valid ed25519 account")
	}

	return nil
}
