package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// The deprecated in-memory API is a thin wrapper around the Ed25519 API, which
// is where the behaviour itself is tested. These tests only check that every
// wrapper delegates to its replacement.

// makeInMemoryTestTxn returns a payment transaction from and to the given address
func makeInMemoryTestTxn(addr types.Address) types.Transaction {
	return types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     addr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: addr,
			Amount:   5000,
		},
	}
}

func TestBasicAccountTransactionSigner(t *testing.T) {
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}
	tx := makeInMemoryTestTxn(account.Address)

	stxs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	expectedSigner := Ed25519AccountTransactionSigner{Signer: account.AsSigner()}
	expectedStxs, err := expectedSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)
	require.Equal(t, expectedStxs, stxs)

	require.True(t, txSigner.Equals(BasicAccountTransactionSigner{Account: account}))
	require.False(t, txSigner.Equals(BasicAccountTransactionSigner{Account: crypto.GenerateAccount()}))
	require.False(t, txSigner.Equals(expectedSigner))
}

func TestMultiSigAccountTransactionSigner(t *testing.T) {
	account1 := crypto.GenerateAccount()
	account2 := crypto.GenerateAccount()
	msig, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{account1.Address, account2.Address})
	require.NoError(t, err)
	msigAddr, err := msig.Address()
	require.NoError(t, err)

	sks := [][]byte{account1.PrivateKey, account2.PrivateKey}
	txSigner := MultiSigAccountTransactionSigner{Msig: msig, Sks: sks}
	tx := makeInMemoryTestTxn(msigAddr)

	stxs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	expectedSigner := MultiSigEd25519AccountTransactionSigner{
		Msig:    msig,
		Signers: []crypto.Ed25519Signer{account1.AsSigner(), account2.AsSigner()},
	}
	expectedStxs, err := expectedSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)
	require.Equal(t, expectedStxs, stxs)

	require.True(t, txSigner.Equals(MultiSigAccountTransactionSigner{Msig: msig, Sks: sks}))
	require.False(t, txSigner.Equals(MultiSigAccountTransactionSigner{Msig: msig, Sks: sks[:1]}))
	require.False(t, txSigner.Equals(expectedSigner))
}
