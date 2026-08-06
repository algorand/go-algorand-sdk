//go:build falcon

package transaction

import (
	"encoding/json"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Falcon1024AccountTransactionSigner is a TransactionSigner that can
// sign transactions for the provided Falcon1024 Account
type Falcon1024AccountTransactionSigner struct {
	Falcon1024Account crypto.Falcon1024Account
}

// SignTransactions signs the provided transactions with the Falcon1024 signer.
func (txSigner Falcon1024AccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		_, stxBytes, err := crypto.SignFalcon1024AccountTransaction(txSigner.Falcon1024Account, txGroup[pos])
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
