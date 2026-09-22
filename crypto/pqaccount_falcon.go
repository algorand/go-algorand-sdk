//go:build falcon

package crypto

import (
	"github.com/algorand/falcon"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Falcon1024PrivateKeySize is the size in bytes of a falcon1024 private key
const Falcon1024PrivateKeySize = 2305

// Falcon1024PublicKeySize is the size in bytes of a falcon1024 public key
const Falcon1024PublicKeySize = 1793

// Falcon1024PublicKey represents a 1793 byte falcon1024 public key.
type Falcon1024PublicKey [Falcon1024PublicKeySize]byte

// Falcon1024Account holds both the public and private information associated with a
// falcon address.
//
// Note: having in-memory cryptographic secrets is discouraged
type Falcon1024Account struct {
	PublicKey  Falcon1024PublicKey
	PrivateKey [Falcon1024PrivateKeySize]byte
}

// Address returns the account address for the given Falcon1024Account, derived
// from its public key with the canonical salt.
func (pqa Falcon1024Account) Address() (addr types.Address, err error) {
	return PQAddress(pqa.PublicKey[:], types.PQSchemeFalcon1024)
}

// basicFalcon1024AccountSigner is a simple signer that wraps an in-memory
// Falcon1024Account
//
// Note: having in-memory cryptographic secrets is discouraged
type basicFalcon1024AccountSigner struct {
	Account Falcon1024Account
}

// PQSign signs the given bytes with a pq signature
func (sgnr basicFalcon1024AccountSigner) PQSign(toBeSigned []byte) ([]byte, error) {
	sk := falcon.PrivateKey(sgnr.Account.PrivateKey)
	return sk.SignCompressed(toBeSigned)
}

// PQPublicKey returns the public key that should be used to verify the
// signatures performed by this signer
func (sgnr basicFalcon1024AccountSigner) PQPublicKey() []byte {
	return sgnr.Account.PublicKey[:]
}

// PQScheme returns the identifier for the post-quantum scheme used by this
// signer
func (sgnr basicFalcon1024AccountSigner) PQScheme() types.PQScheme {
	return types.PQSchemeFalcon1024
}

// AsSigner transforms this account to a PQSigner
//
// Note: having in-memory cryptographic secrets is discouraged
func (pqa Falcon1024Account) AsSigner() PQSigner {
	return &basicFalcon1024AccountSigner{
		Account: pqa,
	}
}

// Falcon1024AccountFromPQSeed returns the corresponding Falcon1024Account for a given seed.
// In conjunction with mnemonic.ToPQSeed() it can be used to generate falcon1024 accounts from regular 25 word mnemonics.
//
// Note: having in-memory cryptographic secrets is discouraged
func Falcon1024AccountFromPQSeed(pqseed []byte) (pqa Falcon1024Account, err error) {
	pk, sk, err := falcon.GenerateKey(pqseed)
	if err != nil {
		return
	}

	pqaPK := Falcon1024PublicKey(pk)
	// fail here rather than at Address() time if this key admits no salt whose
	// address falls outside the ed25519 curve
	if _, err = canonicalSaltForPQPK(pqaPK[:], types.PQSchemeFalcon1024); err != nil {
		return
	}

	pqa = Falcon1024Account{
		PublicKey:  pqaPK,
		PrivateKey: sk,
	}
	return
}

// GenerateFalcon1024Account returns a new Falcon1024Account
//
// Note: having in-memory cryptographic secrets is discouraged
func GenerateFalcon1024Account() Falcon1024Account {
	for {
		seed := make([]byte, 32)
		RandomBytes(seed)
		pka, err := Falcon1024AccountFromPQSeed(seed)
		if err == nil {
			return pka
		}
	}
}

// VerifyPQSig checks that the given pqsig corresponds to the expected toBeSigned byte sequence.
func VerifyPQSig(toBeSigned []byte, pqsig types.PQSig) bool {
	if pqsig.Scheme != types.PQSchemeFalcon1024 {
		return false
	}
	if len(pqsig.PublicKey) != Falcon1024PublicKeySize {
		return false
	}
	pk := falcon.PublicKey(pqsig.PublicKey)
	sig := falcon.CompressedSignature(pqsig.Signature)
	return pk.Verify(sig, toBeSigned) == nil
}

func init() {
	// with falcon available, VerifyLogicSig can cryptographically verify the
	// delegation signature of a PQ-delegated LogicSig instead of only checking
	// the delegating address
	verifyPQDelegation = VerifyPQSig
}
