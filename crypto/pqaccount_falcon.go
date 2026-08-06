//go:build falcon

package crypto

import (
	"fmt"

	"github.com/algorand/falcon"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Falcon1024AccountFromPQSeed returns the corresponding Falcon1024Account for a given seed.
// In conjunction with mnemonic.ToPQSeed() It can be used to generate falcon1024 accounts from regular 25 word mnemonics.
func Falcon1024AccountFromPQSeed(pqseed []byte) (pka Falcon1024Account, err error) {
	pk, sk, err := falcon.GenerateKey(pqseed)
	if err != nil {
		return
	}

	for salt := 0; salt <= 0xff; salt++ {
		pka = Falcon1024Account{
			PublicKey:  pk,
			PrivateKey: sk,
			Salt:       types.PQAddressSalt(salt),
		}

		if pka.Validate() == nil {
			return
		}
	}

	pka = Falcon1024Account{}
	err = fmt.Errorf("no valid salt with an address outside the ed25519 curve exists for the given seed")
	return
}

// GenerateFalcon1024Account returns a new canonical Falcon1024Account
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

// SignFalcon1024AccountTransaction signs the given transaction with the given Falcon1024Account private key. On success it returns both transaction id and transaction bytes.
func SignFalcon1024AccountTransaction(pqa Falcon1024Account, txn types.Transaction) (txid string, stxBytes []byte, err error) {
	txnBytes := rawTransactionBytesToSign(txn)
	txid = txIDFromRawTxnBytesToSign(txnBytes)

	sk := falcon.PrivateKey(pqa.PrivateKey)

	sig, err := sk.SignCompressed(txnBytes[:])
	if err != nil {
		return
	}

	pqsig := types.PQSig{
		Scheme:    types.PQSchemeFalcon1024,
		Salt:      pqa.Salt,
		PublicKey: pqa.PublicKey[:],
		Signature: sig,
	}

	stx := types.SignedTxn{
		Txn:   txn,
		PQsig: pqsig,
	}

	addr := pqa.Address()
	if stx.Txn.Sender != addr {
		stx.AuthAddr = addr
	}

	stxBytes = msgpack.Encode(stx)
	return
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

// MakeLogicSigAccountDelegatedFalcon1024 creates delegated LogicSigAccount that can sign on behalf of a Falcon1024 account.
func MakeLogicSigAccountDelegatedFalcon1024(program []byte, args [][]byte, pqsigAccount Falcon1024Account) (lsa LogicSigAccount, err error) {
	if err = sanityCheckProgram(program); err != nil {
		return
	}

	toSignBytes := pqsigProgramToSign(pqsigAccount.Address(), program)
	sk := falcon.PrivateKey(pqsigAccount.PrivateKey)
	sig, err := sk.SignCompressed(toSignBytes)
	if err != nil {
		return
	}

	pqsig := types.PQSig{
		Scheme:    types.PQSchemeFalcon1024,
		Salt:      pqsigAccount.Salt,
		PublicKey: pqsigAccount.PublicKey[:],
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
