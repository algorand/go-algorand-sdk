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

// SaltedPQSigner wraps a given PQSigner overriding its salt
// with a new one
type SaltedPQSigner struct {
	PQSigner
	Salt types.PQAddressSalt
}

// PQSalt returns the (maybe non-canonical) salt that identifies the
// account selected for this signer
func (sgnr SaltedPQSigner) PQSalt() types.PQAddressSalt {
	return sgnr.Salt
}

// PQAddress returns the account address for the given pq public key, scheme and
// salt.
// Hash("PQA" || scheme || salt || publicKey)
func PQAddress(pk []byte, scheme types.PQScheme, salt types.PQAddressSalt) (addr types.Address) {
	buf := make([]byte, 0, len(pqAddressPrefix)+len(scheme)+1+len(pk))
	buf = append(buf, pqAddressPrefix...)
	buf = append(buf, scheme[:]...)
	buf = append(buf, uint8(salt))
	buf = append(buf, pk[:]...)

	digest := sha512.Sum512_256(buf)

	copy(addr[:], digest[:])
	return
}

// PQSignerAddress returns the address for a given PQSigner
func PQSignerAddress(signer PQSigner) (addr types.Address, err error) {
	salt, err := SaltForPQSigner(signer)
	if err != nil {
		return
	}
	return PQAddress(signer.PQPublicKey(), signer.PQScheme(), salt), nil
}

// SaltForPQSigner returns the salt that will be used when performing PQ
// signatures.
//
// For signers implementing PQSalted this salt will be used, otherwise
// the canonical one will be calculated.
func SaltForPQSigner(sgnr PQSigner) (types.PQAddressSalt, error) {
	if salted, ok := sgnr.(PQSalted); ok {
		return salted.PQSalt(), nil
	}

	return canonicalSaltForPQPK(sgnr.PQPublicKey(), sgnr.PQScheme())
}

func canonicalSaltForPQPK(pk []byte, scheme types.PQScheme) (types.PQAddressSalt, error) {
	for salt := 0; salt <= 0xff; salt++ {
		addr := PQAddress(pk, scheme, types.PQAddressSalt(salt))
		if !IsEdwards25519Point(addr[:]) {
			return types.PQAddressSalt(salt), nil
		}
	}

	return 0, fmt.Errorf("no valid salt with an address outside the ed25519 curve exists for %x", pk)
}

// pqSig assembles the PQSig envelope that identifies the signer's account and
// carries the given signature bytes.
func pqSig(sgnr PQSigner, signature []byte) (sig types.PQSig, addr types.Address, err error) {
	salt, err := SaltForPQSigner(sgnr)
	if err != nil {
		return
	}

	sig = types.PQSig{
		Scheme:    sgnr.PQScheme(),
		Salt:      salt,
		PublicKey: sgnr.PQPublicKey(),
		Signature: signature,
	}
	addr = PQAddress(sig.PublicKey, sig.Scheme, salt)
	return
}

// MakeLogicSigAccountDelegatedPQ creates a delegated LogicSigAccount that can sign on behalf of a PQ account.
func MakeLogicSigAccountDelegatedPQ(program []byte, args [][]byte, sgnr PQSigner) (lsa LogicSigAccount, err error) {
	if err = sanityCheckProgram(program); err != nil {
		return
	}

	// the delegation signature commits to the address it delegates from, so the
	// envelope is assembled first and its signature filled in afterwards
	pqsig, addr, err := pqSig(sgnr, nil)
	if err != nil {
		return
	}

	pqsig.Signature, err = sgnr.PQSign(pqsigProgramToSign(addr, program))
	if err != nil {
		return
	}

	lsa = LogicSigAccount{
		Lsig: types.LogicSig{
			Logic: program,
			Args:  args,
			PQsig: pqsig,
		},
	}
	return
}
