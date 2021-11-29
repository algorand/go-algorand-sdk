package future

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMakeAtomicTransactionComposer(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 0)
	copyAtc := atc.Clone()
	require.Equal(t, atc, copyAtc)
}

func TestAddTransaction(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	tx := types.Transaction{
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

	txAndSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: txSigner,
	}

	err = atc.AddTransaction(txAndSigner)
	require.NoError(t, err)

	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 1)
}

func TestAddTransactionWhenNotBuilding(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	tx := types.Transaction{
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

	txAndSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: txSigner,
	}

	err = atc.AddTransaction(txAndSigner)
	require.NoError(t, err)

	_, err = atc.BuildGroup()
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILT)

	err = atc.AddTransaction(txAndSigner)
	require.Error(t, err)
}

func TestAddTransactionWithMaxTransactions(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	tx := types.Transaction{
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

	txAndSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: txSigner,
	}

	for i := 0; i < 16; i++ {
		err = atc.AddTransaction(txAndSigner)
		require.NoError(t, err)
	}

	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 16)

	err = atc.AddTransaction(txAndSigner)
	require.Error(t, err)
}

func TestAddMethodCall(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}
	methodSig := "add()uint32"

	method, err := MethodFromSignature(methodSig)
	require.NoError(t, err)

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	err = atc.AddMethodCall(
		AddMethodCallParams{
			0,
			method,
			nil,
			addr,
			types.SuggestedParams{},
			types.NoOpOC,
			[]byte{},
			[32]byte{},
			addr,
		},
		txSigner,
	)
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 1)
}

func TestGatherSignatures(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	tx := types.Transaction{
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

	txAndSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: txSigner,
	}

	err = atc.AddTransaction(txAndSigner)
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 1)

	sigs, err := atc.GatherSignatures()
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), SIGNED)
	require.Equal(t, len(sigs), 1)

	tx.Group, err = crypto.ComputeGroupID([]types.Transaction{tx})
	require.NoError(t, err)
	txWithSigners, _ := atc.BuildGroup()
	require.Equal(t, tx.Group, txWithSigners[0].Txn.Group)

	_, expectedSig, err := crypto.SignTransaction(account.PrivateKey, tx)
	require.NoError(t, err)
	require.Equal(t, len(sigs[0]), len(expectedSig))
	require.Equal(t, sigs[0], expectedSig)
}
