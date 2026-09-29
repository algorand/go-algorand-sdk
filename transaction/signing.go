package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/internal/signing"
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
	_, stxBytes, err := signing.Ed25519SignTransaction(signer, tx)
	return stxBytes, err
}

func ed25519SignMultisigTransaction(signer crypto.Ed25519Signer, account crypto.MultisigAccount, tx types.Transaction) ([]byte, error) {
	_, stxBytes, err := signing.Ed25519SignMultisigTransaction(signer, signing.MultisigAccount(account), tx)
	return stxBytes, err
}

func signLogicSigAccountTransaction(account crypto.LogicSigAccount, tx types.Transaction) ([]byte, error) {
	_, stxBytes, err := crypto.SignLogicSigAccountTransaction(account, tx)
	return stxBytes, err
}
