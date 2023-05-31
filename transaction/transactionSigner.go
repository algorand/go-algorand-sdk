package transaction

import (
	"encoding/json"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// TransactionSigner represents a function which can sign transactions from an atomic transaction group.
// @param txnGroup - The atomic group containing transactions to be signed
// @param indexesToSign - An array of indexes in the atomic transaction group that should be signed
// @returns An array of encoded signed transactions. The length of the
//   array will be the same as the length of indexesToSign, and each index i in the array
//   corresponds to the signed transaction from txnGroup[indexesToSign[i]]
type TransactionSigner interface { //nolint:revive // Ignore stuttering for backwards compatibility
	SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error)
	Equals(other TransactionSigner) bool
}

// BasicAccountTransactionSigner that can sign transactions for the provided basic Account.
type BasicAccountTransactionSigner struct {
	Account crypto.Account
}

// SignTransactions signs the provided transactions with the private key of the account.
func (txSigner BasicAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		_, stxBytes, err := crypto.SignTransaction(txSigner.Account.PrivateKey, txGroup[pos])
		if err != nil {
			return nil, err
		}

		stxs[i] = stxBytes
	}

	return stxs, nil
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner BasicAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(BasicAccountTransactionSigner); ok {
		otherJSON, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJSON, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJSON) == string(selfJSON)
	}
	return false
}

// LogicSigAccountTransactionSigner is a TransactionSigner that can
// sign transactions for the provided LogicSigAccount.
type LogicSigAccountTransactionSigner struct {
	LogicSigAccount crypto.LogicSigAccount
}

// SignTransactions signs the provided transactions with the private key of the account.
func (txSigner LogicSigAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		_, stxBytes, err := crypto.SignLogicSigAccountTransaction(txSigner.LogicSigAccount, txGroup[pos])
		if err != nil {
			return nil, err
		}

		stxs[i] = stxBytes
	}

	return stxs, nil
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner LogicSigAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(LogicSigAccountTransactionSigner); ok {
		otherJSON, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJSON, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJSON) == string(selfJSON)
	}
	return false
}

// MultiSigAccountTransactionSigner is a TransactionSigner that can
// sign transactions for the provided MultiSig Account
type MultiSigAccountTransactionSigner struct {
	Msig crypto.MultisigAccount
	Sks  [][]byte
}

// SignTransactions signs the provided transactions with the private keys of the account.
func (txSigner MultiSigAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		var unmergedStxs [][]byte
		for _, sk := range txSigner.Sks {
			_, unmergedStxBytes, err := crypto.SignMultisigTransaction(sk, txSigner.Msig, txGroup[pos])
			if err != nil {
				return nil, err
			}

			unmergedStxs = append(unmergedStxs, unmergedStxBytes)
		}

		if len(txSigner.Sks) > 1 {
			_, stxBytes, err := crypto.MergeMultisigTransactions(unmergedStxs...)
			if err != nil {
				return nil, err
			}

			stxs[i] = stxBytes
		} else {
			stxs[i] = unmergedStxs[0]
		}
	}

	return stxs, nil
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner MultiSigAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(MultiSigAccountTransactionSigner); ok {
		otherJSON, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJSON, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJSON) == string(selfJSON)
	}
	return false
}
