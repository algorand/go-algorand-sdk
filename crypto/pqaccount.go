package crypto

import (
	"crypto/sha512"
	"fmt"

	"github.com/algorand/falcon"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// pqAddressPrefix is prepended when deriving a post-quantum account address.
var pqAddressPrefix = []byte("PQA")

// pqProgramPrefix is prepended to a logic program when computing the bytes a
// post-quantum scheme signs for a delegated LogicSig.
var pqProgramPrefix = []byte("PQProgram")

// Falcon1024Account holds both the public and private information associated with a
// falcon address.
type Falcon1024Account struct {
	PublicKey  falcon.PublicKey
	PrivateKey falcon.PrivateKey
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

// SignFalcon1024AccountTransaction signs the given transaction with the given Falcon1024Account private key. On success it returns both transaction id and transaction bytes.
func SignFalcon1024AccountTransaction(pqa Falcon1024Account, txn types.Transaction) (txid string, stxBytes []byte, err error) {
	txnBytes := rawTransactionBytesToSign(txn)
	txid = txIDFromRawTxnBytesToSign(txnBytes)

	sig, err := pqa.PrivateKey.SignCompressed(txnBytes[:])
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
