package transaction

import (
	"encoding/json"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// TransactionSigner represents a function which can sign transactions from an atomic transaction group.
// @param txnGroup - The atomic group containing transactions to be signed
// @param indexesToSign - An array of indexes in the atomic transaction group that should be signed
// @returns An array of encoded signed transactions. The length of the
//
//	array will be the same as the length of indexesToSign, and each index i in the array
//	corresponds to the signed transaction from txnGroup[indexesToSign[i]]
type TransactionSigner interface { //nolint:revive // Ignore stuttering for backwards compatibility
	SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error)
	Equals(other TransactionSigner) bool
}

// Ed25519AccountTransactionSigner that can sign transactions using the provided Ed25519 signer.
type Ed25519AccountTransactionSigner struct {
	Signer crypto.Ed25519Signer
}

// SignTransactions signs the provided transactions with the private key of the account.
func (txSigner Ed25519AccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		_, stxBytes, err := crypto.Ed25519SignTransaction(txSigner.Signer, txGroup[pos])
		if err != nil {
			return nil, err
		}

		stxs[i] = stxBytes
	}

	return stxs, nil
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner Ed25519AccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(Ed25519AccountTransactionSigner); ok {
		pk1 := txSigner.Signer.Ed25519PublicKey()
		pk2 := castedSigner.Signer.Ed25519PublicKey()
		// NOTE: Assuming that two signers for the same PK are "equal"
		return pk1 == pk2
	}
	return false
}

// MultiSigEd25519AccountTransactionSigner is a TransactionSigner that can sign
// transactions for the provided MultiSig Account
type MultiSigEd25519AccountTransactionSigner struct {
	Msig    crypto.MultisigAccount
	Signers []crypto.Ed25519Signer
}

// SignTransactions signs the provided transactions with the private keys of the account.
func (txSigner MultiSigEd25519AccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		var unmergedStxs [][]byte
		for _, sgnr := range txSigner.Signers {
			_, unmergedStxBytes, err := crypto.Ed25519SignMultisigTransaction(sgnr, txSigner.Msig, txGroup[pos])
			if err != nil {
				return nil, err
			}

			unmergedStxs = append(unmergedStxs, unmergedStxBytes)
		}

		if len(txSigner.Signers) > 1 {
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
func (txSigner MultiSigEd25519AccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(MultiSigEd25519AccountTransactionSigner); ok {
		otherJSON, err := json.Marshal(castedSigner.Msig)
		if err != nil {
			return false
		}

		selfJSON, err := json.Marshal(txSigner.Msig)
		if err != nil {
			return false
		}

		if string(otherJSON) != string(selfJSON) {
			return false
		}

		if len(txSigner.Signers) != len(castedSigner.Signers) {
			return false
		}

		for idx, sgnr := range txSigner.Signers {
			otherSgnr := castedSigner.Signers[idx]
			if sgnr.Ed25519PublicKey() != otherSgnr.Ed25519PublicKey() {
				return false
			}
		}

		return true
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

// Falcon1024AccountTransactionSigner is a TransactionSigner that can
// sign transactions using the provided Falcon1024 signer
type Falcon1024AccountTransactionSigner struct {
	Signer crypto.Falcon1024Signer
}

// SignTransactions signs the provided transactions with the Falcon1024 signer.
func (txSigner Falcon1024AccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		_, stxBytes, err := crypto.SignFalcon1024AccountTransaction(txSigner.Signer, txGroup[pos])
		if err != nil {
			return nil, err
		}

		stxs[i] = stxBytes
	}

	return stxs, nil
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner Falcon1024AccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(Falcon1024AccountTransactionSigner); ok {
		// NOTE: Assuming that two signers for the same (PK, salt) pair are "equal"
		if txSigner.Signer.Falcon1024PublicKey() != castedSigner.Signer.Falcon1024PublicKey() {
			return false
		}
		txSignerSalt, txSignerSaltErr := crypto.SaltForFalcon1024Signer(txSigner.Signer)
		castedSignerSalt, castedSignerSaltErr := crypto.SaltForFalcon1024Signer(castedSigner.Signer)
		// If both fail then they are the same weird thing
		if txSignerSaltErr != castedSignerSaltErr {
			return false
		}
		// Otherwise they should have the same salt
		return txSignerSalt == castedSignerSalt
	}
	return false
}

// EmptyTransactionSigner is a TransactionSigner that produces signed transaction objects without
// signatures. This is useful for simulating transactions, but it won't work for actual submission.
type EmptyTransactionSigner struct{}

// SignTransactions returns SignedTxn bytes but does not sign them.
func (txSigner EmptyTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		stx := types.SignedTxn{
			Txn: txGroup[pos],
		}
		stxs[i] = msgpack.Encode(&stx)
	}
	return stxs, nil
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner EmptyTransactionSigner) Equals(other TransactionSigner) bool {
	_, ok := other.(EmptyTransactionSigner)
	return ok
}
