//go:build falcon

package crypto

import (
	"fmt"

	"github.com/algorand/falcon"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Falcon1024PrivateKeySize is the size in bytes of a falcon1024 private key
const Falcon1024PrivateKeySize = 2305

// Falcon1024Account holds both the public and private information associated with a
// falcon address.
//
// Note: having in-memory cryptographic secrets is discouraged
type Falcon1024Account struct {
	PublicKey  Falcon1024PublicKey
	PrivateKey [Falcon1024PrivateKeySize]byte
	Salt       types.PQAddressSalt
}

// Address returns the account address for the given Falcon1024Account.
// Hash("PQA" || scheme || salt || publicKey)
func (pqa Falcon1024Account) Address() (addr types.Address) {
	return falcon1024Address(pqa.PublicKey, pqa.Salt)
}

// Validate returns an error if the given Falcon1024Account address could be interpreted as an ed25519 public key
func (pqa Falcon1024Account) Validate() error {
	addr := pqa.Address()

	if IsEdwards25519Point([]byte(addr[:])) {
		return fmt.Errorf("Account address overlaps with a valid ed25519 account")
	}

	return nil
}

// basicFalcon1024AccountSigner is a simple signer that wraps an in-memory
// Falcon1024Account
//
// Note: having in-memory cryptographic secrets is discouraged
type basicFalcon1024AccountSigner struct {
	Account Falcon1024Account
}

// Falcon1024Sign signs the given bytes with a falcon1024 signature
func (sgnr basicFalcon1024AccountSigner) Falcon1024Sign(toBeSigned []byte) ([]byte, error) {
	sk := falcon.PrivateKey(sgnr.Account.PrivateKey)
	return sk.SignCompressed(toBeSigned)
}

// Falcon1024PublicKey returns the public key that should be used to verify the
// signatures performed by this signer
func (sgnr basicFalcon1024AccountSigner) Falcon1024PublicKey() Falcon1024PublicKey {
	return sgnr.Account.PublicKey
}

// Falcon1024Salt returns the (maybe non-canonical) salt that identifies the
// account selected for this signer
func (sgnr basicFalcon1024AccountSigner) Falcon1024Salt() types.PQAddressSalt {
	return sgnr.Account.Salt
}

// AsSigner transforms this account to a Falcon1024Signer
//
// The resulting signer will respect the salt of the source account.
//
// Note: having in-memory cryptographic secrets is discouraged
func (pqa Falcon1024Account) AsSigner() Falcon1024Signer {
	return basicFalcon1024AccountSigner{
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
	salt, err := canonicalSaltForFalcon1024PK(pqaPK)
	if err != nil {
		return
	}

	pqa = Falcon1024Account{
		PublicKey:  pqaPK,
		PrivateKey: sk,
		Salt:       salt,
	}
	return
}

// GenerateFalcon1024Account returns a new canonical Falcon1024Account
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
	pk := falcon.PublicKey(pqsig.PublicKey)
	sig := falcon.CompressedSignature(pqsig.Signature)
	return pk.Verify(sig, toBeSigned) == nil
}

// verifyPQSig is the unexported version of the VerifyPQSig used by VerifyLogicSig.
// On falcon-disabled environment the implementation will refuse to validate the signature.
// See pqaccount_nofalcon.go for more info.
func verifyPQSig(toBeSigned []byte, pqsig types.PQSig) bool {
	return VerifyPQSig(toBeSigned, pqsig)
}
