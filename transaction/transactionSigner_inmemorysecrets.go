package transaction

import (
	"encoding/json"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// BasicAccountTransactionSigner that can sign transactions for the provided basic Account.
//
// Deprecated: having in-memory cryptographic secrets is discouraged, use
// Ed25519AccountTransactionSigner instead
type BasicAccountTransactionSigner struct {
	Account crypto.Account
}

// SignTransactions signs the provided transactions with the private key of the account.
func (txSigner BasicAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	transactionSigner := Ed25519AccountTransactionSigner{Signer: txSigner.Account.AsSigner()}
	return transactionSigner.SignTransactions(txGroup, indexesToSign)
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

// MultiSigAccountTransactionSigner is a TransactionSigner that can
// sign transactions for the provided MultiSig Account
//
// Deprecated: having in-memory cryptographic secrets is discouraged, use
// MultiSigEd25519AccountTransactionSigner instead
type MultiSigAccountTransactionSigner struct {
	Msig crypto.MultisigAccount
	Sks  [][]byte
}

// SignTransactions signs the provided transactions with the private keys of the account.
func (txSigner MultiSigAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	signers := make([]crypto.Ed25519Signer, len(txSigner.Sks))
	for i, sk := range txSigner.Sks {
		signer, err := crypto.SKToInMemorySigner(sk)
		if err != nil {
			return nil, err
		}

		signers[i] = signer
	}
	transactionSigner := MultiSigEd25519AccountTransactionSigner{Msig: txSigner.Msig, Signers: signers}
	return transactionSigner.SignTransactions(txGroup, indexesToSign)
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
