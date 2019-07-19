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

// VerifySignature checks that the given hashed data is has a valid signature.
func VerifySignature(pk ed25519.PublicKey, data []byte, sig types.Signature) bool {
	return ed25519.Verify(pk, data, sig.ToBytes())
}

// SignTransactionRaw returns an unencoded version of the transaction.
func SignTransactionRaw(sk ed25519.PrivateKey, tx types.Transaction) (s types.Signature, txid string, stx types.SignedTxn, err error) {
  s, txid, err = rawSignTransaction(sk, tx)
  if err != nil {
    return
  }
  // Construct the SignedTxn
  stx = types.SignedTxn{
    Sig: s,
    Txn: tx,
  }
	return
}

// SignTransaction accepts a private key and a transaction, and returns the
// bytes of a signed transaction ready to be broadcasted to the network
func SignTransaction(sk ed25519.PrivateKey, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	s, txid, err := rawSignTransaction(sk, tx)
	if err != nil {
		return
	}
	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

// rawTransactionBytesToSign returns the byte form of the tx that we actually sign
// and compute txID from.
func rawTransactionBytesToSign(tx types.Transaction) ([]byte) {
	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Prepend the hashable prefix
	msgParts := [][]byte{txidPrefix, encodedTx}
	return bytes.Join(msgParts, nil)
}

// txID computes a transaction id from raw transaction bytes
func txIDFromRawTxnBytesToSign(toBeSigned []byte) (txid string) {
	txidBytes := sha512.Sum512_256(toBeSigned)
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
	return
}

// txIDFromTransaction is a convenience function for generating txID from txn
func txIDFromTransaction(tx types.Transaction) (txid string) {
	txid = txIDFromRawTxnBytesToSign(rawTransactionBytesToSign(tx))
	return
}

// rawSignTransaction signs the msgpack-encoded tx (with prepended "TX" prefix), and returns the sig and txid
func rawSignTransaction(sk ed25519.PrivateKey, tx types.Transaction) (s types.Signature, txid string, err error) {
	toBeSigned := rawTransactionBytesToSign(tx)

	// Sign the encoded transaction
	signature := ed25519.Sign(sk, toBeSigned)

	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	n := copy(s[:], signature)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}
	// Populate txID
	txid = txIDFromRawTxnBytesToSign(toBeSigned)
	return
}

// SignBid accepts a private key and a bid, and returns the signature of the
// bid under that key
func SignBid(sk ed25519.PrivateKey, bid types.Bid) (signedBid []byte, err error) {
	// Encode the bid as msgpack
	encodedBid := msgpack.Encode(bid)

	// Prepend the hashable prefix
	msgParts := [][]byte{bidPrefix, encodedBid}
	toBeSigned := bytes.Join(msgParts, nil)

	// Sign the encoded bid
	sig := ed25519.Sign(sk, toBeSigned)

	var s types.Signature
	n := copy(s[:], sig)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}

	sb := types.SignedBid{
		Bid: bid,
		Sig: s,
	}

	nf := types.NoteField{
		Type:      types.NoteBid,
		SignedBid: sb,
	}

	signedBid = msgpack.Encode(nf)
	return
}

/* Multisig Support */

// SignMultisigTransaction signs the given transaction, and multisig preimage, with the
// private key, returning the bytes of a signed transaction with the multisig field
// partially populated, ready to be passed to other multisig signers to sign or broadcast.
func SignMultisigTransaction(sk ed25519.PrivateKey, pk MultisigAccount, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	err = pk.Validate()
	if err != nil {
		return
	}
	// check that the address of txn matches the preimage
	pkAddr, err := pk.Address()
	if err != nil {
		return
	}
	if tx.Sender != pkAddr { // array value comparison is fine
		err = errMsigBadTxnSender
		return
	}
	// check that sk.pk exists in the pk list
	myIndex := len(pk.Pks)
	myPk := sk.Public().(ed25519.PublicKey)
	for i := 0; i < len(pk.Pks); i++ {
		if bytes.Equal(myPk, pk.Pks[i]) {
			myIndex = i
		}
	}
	if myIndex == len(pk.Pks) {
		err = errMsigInvalidSecretKey
		return
	}
	// now, create the signed transaction
	var sig types.MultisigSig
	sig.Version = pk.Version
	sig.Threshold = pk.Threshold
	sig.Subsigs = make([]types.MultisigSubsig, len(pk.Pks))
	for i := 0; i < len(pk.Pks); i++ {
		c := make([]byte, len(pk.Pks[i]))
		copy(c, pk.Pks[i])
		sig.Subsigs[i].Key = c
	}
	rawSig, txid, err := rawSignTransaction(sk, tx)
	if err != nil {
		return
	}
	sig.Subsigs[myIndex].Sig = rawSig

	// Encode the signedTxn
	stx := types.SignedTxn{
		Msig: sig,
		Txn: tx,
	}
	stxBytes = msgpack.Encode(stx)
	return
}

// MergeMultisigTransactions merges the given (partially) signed multisig transactions, and
// returns an encoded signed multisig transaction with the component signatures.
func MergeMultisigTransactions(stxsBytes ...[]byte) (txid string, stxBytes []byte, err error) {
	if len(stxsBytes) < 2 {
		err = errMsigMergeLessThanTwo
		return
	}
	var sig types.MultisigSig
	var refAddr *types.Address
	var refTx types.Transaction
	for _, partStxBytes := range stxsBytes {
		partStx := types.SignedTxn{}
		err = msgpack.Decode(partStxBytes, &partStx)
		if err != nil {
			return
		}
		// check that multisig parameters match
		partMa, innerErr := MultisigAccountFromSig(partStx.Msig)
		if innerErr != nil {
			err = innerErr
			return
		}
		partAddr, innerErr := partMa.Address()
		if innerErr != nil {
			err = innerErr
			return
		}
		if refAddr == nil {
			refAddr = &partAddr
			// add parameters to new merged txn
			sig.Version = partStx.Msig.Version
			sig.Threshold = partStx.Msig.Threshold
			sig.Subsigs = make([]types.MultisigSubsig, len(partStx.Msig.Subsigs))
			for i := 0; i < len(sig.Subsigs); i++ {
				c := make([]byte, len(partStx.Msig.Subsigs[i].Key))
				copy(c, partStx.Msig.Subsigs[i].Key)
				sig.Subsigs[i].Key = c
			}
			refTx = partStx.Txn
		} else {
			if partAddr != *refAddr {
				err = errMsigMergeKeysMismatch
				return
			}
		}
		// now, add subsignatures appropriately
		zeroSig := types.Signature{}
		for i := 0; i < len(sig.Subsigs); i++ {
			mSubsig := partStx.Msig.Subsigs[i]
			if mSubsig.Sig != zeroSig {
				if sig.Subsigs[i].Sig == zeroSig {
					sig.Subsigs[i].Sig = mSubsig.Sig
				} else if sig.Subsigs[i].Sig != mSubsig.Sig {
					err = errMsigMergeInvalidDups
					return
				}
			}
		}
	}
	// Encode the signedTxn
	stx := types.SignedTxn{
		Msig: sig,
		Txn: refTx,
	}
	stxBytes = msgpack.Encode(stx)
	// let's also compute the txid.
	txid = txIDFromTransaction(refTx)
	return
}

// AppendMultisigTransaction appends the signature corresponding to the given private key,
// returning an encoded signed multisig transaction including the signature.
// While we could compute the multisig preimage from the multisig blob, we ask the caller
// to pass it back in, to explicitly check that they know who they are signing as.
func AppendMultisigTransaction(sk ed25519.PrivateKey, pk MultisigAccount, preStxBytes []byte) (txid string, stxBytes []byte, err error) {
	preStx := types.SignedTxn{}
	err = msgpack.Decode(preStxBytes, &preStx)
	if err != nil {
		return
	}
	_, partStxBytes, err := SignMultisigTransaction(sk, pk, preStx.Txn)
	if err != nil {
		return
	}
	txid, stxBytes, err = MergeMultisigTransactions(partStxBytes, preStxBytes)
	return
}

func GenerateKeysFromSeed(seed []byte) (*ed25519.PublicKey, *ed25519.PrivateKey, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, nil, fmt.Errorf("seed from length mismatch: %d != %d", len(seed), ed25519.SeedSize)
	}

	var publicKey ed25519.PublicKey = make([]byte, ed25519.PublicKeySize)
	privateKey := ed25519.NewKeyFromSeed(seed)
	copy(publicKey, privateKey[32:])
	return &publicKey, &privateKey, nil
}
