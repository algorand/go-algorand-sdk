package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// RandomBytes fills the passed slice with randomness, and panics if it is
// unable to do so
func RandomBytes(s []byte) {
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
}

// SignTransaction accepts a private key and a transaction, and returns the
// bytes of a signed transaction ready to be broadcasted to the network
func SignTransaction(sk ed25519.PrivateKey, tx types.Transaction) ([]byte, error) {
	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Sign the encoded transaction
	signature, err := sk.Sign(nil, encodedTx, nil)
	if err != nil {
		return nil, err
	}

	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	var s types.Signature
	n := copy(s[:], signature)
	if n != len(s) {
		return nil, errInvalidSignatureReturned
	}

	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	// Encode the SignedTxn
	return msgpack.Encode(stx), nil
}
