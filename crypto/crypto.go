package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// txidPrefix is prepended to a transaction when computing its txid
var txidPrefix = []byte("TX")

// bidPrefix is prepended to a bid when signing it
var bidPrefix = []byte("aB")

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
func SignTransaction(sk ed25519.PrivateKey, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Prepend the hashable prefix
	msgParts := [][]byte{txidPrefix, encodedTx}
	toBeSigned := bytes.Join(msgParts, nil)

	// Sign the encoded transaction
	signature := ed25519.Sign(sk, toBeSigned)

	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	var s types.Signature
	n := copy(s[:], signature)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}

	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)

	// Compute the txid
	txidBytes := sha512.Sum512_256(toBeSigned)
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
	return
}

// SignBid accepts a private key and a bid, and returns the signature of the
// bid under that key
func SignBid(sk ed25519.PrivateKey, bid types.Bid) (sig []byte, err error) {
	// Encode the bid as msgpack
	encodedBid := msgpack.Encode(bid)

	// Prepend the hashable prefix
	msgParts := [][]byte{bidPrefix, encodedBid}
	toBeSigned := bytes.Join(msgParts, nil)

	// Sign the encoded bid
	sig = ed25519.Sign(sk, toBeSigned)
	return
}
