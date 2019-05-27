package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"fmt"

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
// bytes of a signed bid in a note.
func SignTransaction(sk []byte, encodedTx []byte) (stxBytes []byte, err error) {
	if len(sk) != ed25519.PrivateKeySize {
		err = fmt.Errorf("Incorrect pricateKey length expected %d, got %d", ed25519.PrivateKeySize, len(sk))
		return
	}

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

	var tx types.Transaction
	msgpack.Decode(encodedTx, &tx)

	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

// SignBid accepts a private key and a bid, and returns the signature of the
// bid under that key
func SignBid(sk []byte, encodedBid []byte) (sBid []byte, err error) {
	if len(sk) != ed25519.PrivateKeySize {
		err = fmt.Errorf("Incorrect pricateKey length expected %d, got %d", ed25519.PrivateKeySize, len(sk))
		return
	}

	// Prepend the hashable prefix
	msgParts := [][]byte{bidPrefix, encodedBid}
	toBeSigned := bytes.Join(msgParts, nil)

	// Sign the encoded bid
	signature := ed25519.Sign(sk, toBeSigned)

	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	var s types.Signature
	n := copy(s[:], signature)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}

	var bid types.Bid
	err = msgpack.Decode(encodedBid, &bid)
	if err != nil {
		return
	}

	//Construct Signed Bid

	signedBid := types.SignedBid{
		Bid: bid,
		Sig: s,
	}

	note := types.NoteField{
		Type:      types.NoteBid,
		SignedBid: signedBid,
	}

	sBid = msgpack.Encode(note)
	return
}

// GetTxID takes an encoded txn and return the txid as string
func GetTxID(encodedTxn []byte) string {
	// Prepend the hashable prefix
	msgParts := [][]byte{txidPrefix, encodedTxn}
	toBeSigned := bytes.Join(msgParts, nil)

	// Compute the txid
	txidBytes := sha512.Sum512_256(toBeSigned)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])

}
