package transaction

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/stretchr/testify/require"
)

func TestGatherSignaturesWithNilSignerReturnsError(t *testing.T) {
	var composer AtomicTransactionComposer
	require.NoError(t, composer.AddTransaction(TransactionWithSigner{
		Txn: types.Transaction{},
	}))

	_, err := composer.GatherSignatures()
	require.ErrorContains(t, err, "nil signer")
}

func TestSignTransactionWithNilSignerReturnsError(t *testing.T) {
	_, _, err := SignTransaction(nil, types.Transaction{})
	require.ErrorContains(t, err, "nil signer")
}
