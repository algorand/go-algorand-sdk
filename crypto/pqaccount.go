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

// SaltedFalcon1024Signer wraps a given Falcon1024Signer overriding its salt
// with a new one
type SaltedFalcon1024Signer struct {
	Signer Falcon1024Signer
	Salt   types.PQAddressSalt
}

// Falcon1024Sign signs the given bytes with a falcon1024 signature
func (sgnr SaltedFalcon1024Signer) Falcon1024Sign(toBeSigned []byte) ([]byte, error) {
	return sgnr.Signer.Falcon1024Sign(toBeSigned)
}

// Falcon1024PublicKey returns the public key that should be used to verify the
// signatures performed by this signer
func (sgnr SaltedFalcon1024Signer) Falcon1024PublicKey() Falcon1024PublicKey {
	return sgnr.Signer.Falcon1024PublicKey()
}

// Falcon1024Salt returns the (maybe non-canonical) salt that identifies the
// account selected for this signer
func (sgnr SaltedFalcon1024Signer) Falcon1024Salt() types.PQAddressSalt {
	return sgnr.Salt
}

// falcon1024Address returns the account address for the given Falcon1024Account.
// Hash("PQA"  || scheme ||  salt || publicKey)
func falcon1024Address(pk Falcon1024PublicKey, salt types.PQAddressSalt) (addr types.Address) {
	buf := make([]byte, 0, len(pqAddressPrefix)+len(types.PQSchemeFalcon1024)+1+len(pk))
	buf = append(buf, pqAddressPrefix...)
	buf = append(buf, types.PQSchemeFalcon1024[:]...)
	buf = append(buf, uint8(salt))
	buf = append(buf, pk[:]...)

	digest := sha512.Sum512_256(buf)

	copy(addr[:], digest[:])
	return
}

// Falcon1024SignerAddress returns the address for a given falcon1024 signer
func Falcon1024SignerAddress(signer Falcon1024Signer) (addr types.Address, err error) {
	salt, err := SaltForFalcon1024Signer(signer)
	if err != nil {
		return
	}
	return falcon1024Address(signer.Falcon1024PublicKey(), salt), nil
}

// SaltForFalcon1024Signer returns the salt that will be used when performing PQ
// signatures.
//
// For signers implementing Falcon1024Salted this salt will be used, otherwise
// the canonical one will be calculated.
func SaltForFalcon1024Signer(sgnr Falcon1024Signer) (types.PQAddressSalt, error) {
	if salted, ok := sgnr.(Falcon1024Salted); ok {
		return salted.Falcon1024Salt(), nil
	}

	return canonicalSaltForFalcon1024PK(sgnr.Falcon1024PublicKey())
}

func canonicalSaltForFalcon1024PK(pk Falcon1024PublicKey) (types.PQAddressSalt, error) {
	for salt := 0; salt <= 0xff; salt++ {
		addr := falcon1024Address(pk, types.PQAddressSalt(salt))
		if !IsEdwards25519Point(addr[:]) {
			return types.PQAddressSalt(salt), nil
		}
	}

	return 0, fmt.Errorf("no valid salt with an address outside the ed25519 curve exists for %s", pk)
}

// SignFalcon1024AccountTransaction signs the given transaction with the given Falcon1024Account private key. On success it returns both transaction id and transaction bytes.
func SignFalcon1024AccountTransaction(sgnr Falcon1024Signer, txn types.Transaction) (txid string, stxBytes []byte, err error) {
	txnBytes := rawTransactionBytesToSign(txn)
	txid = txIDFromRawTxnBytesToSign(txnBytes)

	sig, err := sgnr.Falcon1024Sign(txnBytes)
	if err != nil {
		return
	}

	pk := sgnr.Falcon1024PublicKey()
	salt, err := SaltForFalcon1024Signer(sgnr)
	if err != nil {
		return
	}

	pqsig := types.PQSig{
		Scheme:    types.PQSchemeFalcon1024,
		Salt:      salt,
		PublicKey: pk[:],
		Signature: sig,
	}

	stx := types.SignedTxn{
		Txn:   txn,
		PQsig: pqsig,
	}

	addr := falcon1024Address(pk, salt)
	if stx.Txn.Sender != addr {
		stx.AuthAddr = addr
	}

	stxBytes = msgpack.Encode(stx)
	return
}

// MakeLogicSigAccountDelegatedFalcon1024 creates delegated LogicSigAccount that can sign on behalf of a Falcon1024 account.
func MakeLogicSigAccountDelegatedFalcon1024(program []byte, args [][]byte, sgnr Falcon1024Signer) (lsa LogicSigAccount, err error) {
	if err = sanityCheckProgram(program); err != nil {
		return
	}

	pk := sgnr.Falcon1024PublicKey()
	salt, err := SaltForFalcon1024Signer(sgnr)
	if err != nil {
		return
	}

	addr := falcon1024Address(pk, salt)
	toSignBytes := pqsigProgramToSign(addr, program)
	sig, err := sgnr.Falcon1024Sign(toSignBytes)
	if err != nil {
		return
	}

	pqsig := types.PQSig{
		Scheme:    types.PQSchemeFalcon1024,
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
