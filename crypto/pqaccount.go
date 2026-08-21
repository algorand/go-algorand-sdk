package crypto

import (
	"crypto/sha512"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
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
	Signer PQSigner
	Salt   types.PQAddressSalt
}

// PQSign signs the given bytes with a pq signature
func (sgnr SaltedPQSigner) PQSign(toBeSigned []byte) ([]byte, error) {
	return sgnr.Signer.PQSign(toBeSigned)
}

// PQPublicKey returns the public key that should be used to verify the
// signatures performed by this signer
func (sgnr SaltedPQSigner) PQPublicKey() []byte {
	return sgnr.Signer.PQPublicKey()
}

// PQScheme returns the identifier for the post-quantum scheme used by this
// signer
func (sgnr SaltedPQSigner) PQScheme() types.PQScheme {
	return sgnr.Signer.PQScheme()
}

// PQSalt returns the (maybe non-canonical) salt that identifies the
// account selected for this signer
func (sgnr SaltedPQSigner) PQSalt() types.PQAddressSalt {
	return sgnr.Salt
}

// pqAddress returns the account address for the given pq public key, scheme and
// salt.
// Hash("PQA" || scheme || salt || publicKey)
func pqAddress(pk []byte, scheme types.PQScheme, salt types.PQAddressSalt) (addr types.Address) {
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
	return pqAddress(signer.PQPublicKey(), signer.PQScheme(), salt), nil
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
		addr := pqAddress(pk, scheme, types.PQAddressSalt(salt))
		if !IsEdwards25519Point(addr[:]) {
			return types.PQAddressSalt(salt), nil
		}
	}

	return 0, fmt.Errorf("no valid salt with an address outside the ed25519 curve exists for %s", pk)
}

// SignPQAccountTransaction signs the given transaction with the given PQSigner. On success it returns both transaction id and transaction bytes.
func SignPQAccountTransaction(sgnr PQSigner, txn types.Transaction) (txid string, stxBytes []byte, err error) {
	txnBytes := rawTransactionBytesToSign(txn)
	txid = txIDFromRawTxnBytesToSign(txnBytes)

	sig, err := sgnr.PQSign(txnBytes)
	if err != nil {
		return
	}

	pk := sgnr.PQPublicKey()
	salt, err := SaltForPQSigner(sgnr)
	if err != nil {
		return
	}

	pqsig := types.PQSig{
		Scheme:    sgnr.PQScheme(),
		Salt:      salt,
		PublicKey: pk[:],
		Signature: sig,
	}

	stx := types.SignedTxn{
		Txn:   txn,
		PQsig: pqsig,
	}

	addr := pqAddress(pk, sgnr.PQScheme(), salt)
	if stx.Txn.Sender != addr {
		stx.AuthAddr = addr
	}

	stxBytes = msgpack.Encode(stx)
	return
}

// MakeLogicSigAccountDelegatedPQ creates a delegated LogicSigAccount that can sign on behalf of a PQ account.
func MakeLogicSigAccountDelegatedPQ(program []byte, args [][]byte, sgnr PQSigner) (lsa LogicSigAccount, err error) {
	if err = sanityCheckProgram(program); err != nil {
		return
	}

	pk := sgnr.PQPublicKey()
	salt, err := SaltForPQSigner(sgnr)
	if err != nil {
		return
	}

	addr := pqAddress(pk, sgnr.PQScheme(), salt)
	toSignBytes := pqsigProgramToSign(addr, program)
	sig, err := sgnr.PQSign(toSignBytes)
	if err != nil {
		return
	}

	pqsig := types.PQSig{
		Scheme:    sgnr.PQScheme(),
		Salt:      salt,
		PublicKey: pk[:],
		Signature: sig,
	}

	lsig := types.LogicSig{
		Logic: program,
		Args:  args,
		PQsig: pqsig,
	}

	lsa = LogicSigAccount{
		Lsig: lsig,
	}
	return
}
