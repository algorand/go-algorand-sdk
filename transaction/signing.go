package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// SignTransaction signs one transaction with the provided TransactionSigner.
func SignTransaction(signer TransactionSigner, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	stxs, err := signer.SignTransactions([]types.Transaction{tx}, []int{0})
	if err != nil {
		return "", nil, err
	}
	if len(stxs) != 1 {
		return "", nil, fmt.Errorf("transaction signer returned %d signed transactions, expected 1", len(stxs))
	}
	return crypto.GetTxID(tx), stxs[0], nil
}

// signTransactions signs the transactions of the group selected by
// indexesToSign with signOne, returning the encoded signed transactions in the
// same order as indexesToSign.
func signTransactions(txGroup []types.Transaction, indexesToSign []int, signOne func(types.Transaction) ([]byte, error)) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		stxBytes, err := signOne(txGroup[pos])
		if err != nil {
			return nil, err
		}

		stxs[i] = stxBytes
	}

	return stxs, nil
}

// encodeSignedTxn encodes the SignedTxn, assigning the signing account as the
// AuthAddr when it is not the sender of the transaction.
func encodeSignedTxn(stx types.SignedTxn, signerAddress types.Address) []byte {
	if stx.Txn.Sender != signerAddress {
		stx.AuthAddr = signerAddress
	}
	return msgpack.Encode(&stx)
}

// equalBySerialization reports whether two signers hold equivalent parameters,
// by comparing their JSON encodings.
func equalBySerialization(signer, other interface{}) bool {
	signerJSON, err := json.Marshal(signer)
	if err != nil {
		return false
	}

	otherJSON, err := json.Marshal(other)
	if err != nil {
		return false
	}

	return bytes.Equal(signerJSON, otherJSON)
}

// equalSignerImplementations reports whether two signer interface values refer
// to the same implementation. In particular, pointers are compared by identity
// rather than by the account they sign for.
func equalSignerImplementations(signer, other interface{}) bool {
	if signer == nil || other == nil {
		return signer == nil && other == nil
	}

	signerType := reflect.TypeOf(signer)
	if signerType != reflect.TypeOf(other) || !signerType.Comparable() {
		return false
	}
	return signer == other
}

func ed25519SignTransaction(signer crypto.Ed25519Signer, tx types.Transaction) ([]byte, error) {
	sig, err := crypto.Ed25519RawSignature(signer, crypto.TransactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}

	return encodeSignedTxn(types.SignedTxn{Sig: sig, Txn: tx}, types.Address(signer.Ed25519PublicKey())), nil
}

func ed25519SignMultisigTransaction(signer crypto.Ed25519Signer, account crypto.MultisigAccount, tx types.Transaction) ([]byte, error) {
	if err := account.Validate(); err != nil {
		return nil, err
	}

	msig, err := crypto.Ed25519MultisigSigWith(signer, account, crypto.TransactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}

	address, err := account.Address()
	if err != nil {
		return nil, err
	}
	return encodeSignedTxn(types.SignedTxn{Msig: msig, Txn: tx}, address), nil
}

func ed25519AppendMultisigTransaction(signer crypto.Ed25519Signer, account crypto.MultisigAccount, encoded []byte) (txid string, stxBytes []byte, err error) {
	var stx types.SignedTxn
	if err = msgpack.Decode(encoded, &stx); err != nil {
		return
	}
	partial, err := ed25519SignMultisigTransaction(signer, account, stx.Txn)
	if err != nil {
		return "", nil, err
	}
	return crypto.MergeMultisigTransactions(partial, encoded)
}

func signLogicSigAccountTransaction(account crypto.LogicSigAccount, tx types.Transaction) ([]byte, error) {
	_, stxBytes, err := crypto.SignLogicSigAccountTransaction(account, tx)
	return stxBytes, err
}

// pqSignedTxn returns the encoded SignedTxn carrying the given post-quantum
// signature bytes on behalf of the signer's account.
//
// The signature may be empty, which produces an envelope suitable for
// simulating transactions with the allowEmptySignatures option enabled.
func pqSignedTxn(signer crypto.PQSigner, tx types.Transaction, signature []byte) ([]byte, error) {
	pqsig, address, err := crypto.PQSigFor(signer, signature)
	if err != nil {
		return nil, err
	}

	return encodeSignedTxn(types.SignedTxn{Txn: tx, PQsig: pqsig}, address), nil
}

func pqSignTransaction(signer crypto.PQSigner, tx types.Transaction) ([]byte, error) {
	if signer == nil {
		return nil, crypto.ErrNilPQSigner
	}

	signature, err := signer.PQSign(crypto.TransactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}

	return pqSignedTxn(signer, tx, signature)
}
