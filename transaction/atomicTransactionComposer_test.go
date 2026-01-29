package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func TestMakeAtomicTransactionComposer(t *testing.T) {
	var atc AtomicTransactionComposer
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 0)
	copyAtc := atc.Clone()
	require.Equal(t, atc, copyAtc)
}

func TestAddTransaction(t *testing.T) {
	var atc AtomicTransactionComposer
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
	var atc AtomicTransactionComposer
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
	var atc AtomicTransactionComposer
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
	var atc AtomicTransactionComposer
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}
	methodSig := "add()uint32"

	method, err := abi.MethodFromSignature(methodSig)
	require.NoError(t, err)

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	err = atc.AddMethodCall(
		AddMethodCallParams{
			AppID:  4,
			Method: method,
			Sender: addr,
			Signer: txSigner,
		})
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 1)
}

func TestAddMethodCallWithManualForeignArgs(t *testing.T) {
	var atc AtomicTransactionComposer
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}
	methodSig := "add(application)uint32"

	method, err := abi.MethodFromSignature(methodSig)
	require.NoError(t, err)

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	arg_addr_str := "E4VCHISDQPLIZWMALIGNPK2B2TERPDMR64MZJXE3UL75MUDXZMADX5OWXM"
	arg_addr, err := types.DecodeAddress(arg_addr_str)
	require.NoError(t, err)

	params := AddMethodCallParams{
		AppID:           4,
		Method:          method,
		Sender:          addr,
		Signer:          txSigner,
		MethodArgs:      []interface{}{2},
		ForeignApps:     []uint64{1},
		ForeignAssets:   []uint64{5},
		ForeignAccounts: []string{arg_addr_str},
	}
	err = atc.AddMethodCall(params)
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 1)
	txns, err := atc.BuildGroup()
	require.NoError(t, err)

	require.Equal(t, len(txns[0].Txn.ForeignApps), 2)
	require.Equal(t, txns[0].Txn.ForeignApps[0], types.AppIndex(1))
	require.Equal(t, txns[0].Txn.ForeignApps[1], types.AppIndex(2))
	// verify original params object hasn't changed.
	require.Equal(t, params.ForeignApps, []uint64{1})

	require.Equal(t, len(txns[0].Txn.ForeignAssets), 1)
	require.Equal(t, txns[0].Txn.ForeignAssets[0], types.AssetIndex(5))

	require.Equal(t, len(txns[0].Txn.Accounts), 1)
	require.Equal(t, txns[0].Txn.Accounts[0], arg_addr)

	// make sure Access field is not used with ref args
	params.UseAccess = true
	err = atc.AddMethodCall(params)
	require.Error(t, err)
}

func TestAddMethodCallWithAccess(t *testing.T) {
	var atc AtomicTransactionComposer
	account := crypto.GenerateAccount()
	txSigner := BasicAccountTransactionSigner{Account: account}
	methodSig := "add()uint32"

	method, err := abi.MethodFromSignature(methodSig)
	require.NoError(t, err)

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	arg_addr_str := "E4VCHISDQPLIZWMALIGNPK2B2TERPDMR64MZJXE3UL75MUDXZMADX5OWXM"
	arg_addr, err := types.DecodeAddress(arg_addr_str)
	require.NoError(t, err)

	params := AddMethodCallParams{
		AppID:           4,
		Method:          method,
		Sender:          addr,
		Signer:          txSigner,
		ForeignApps:     []uint64{1},
		ForeignAssets:   []uint64{5},
		ForeignAccounts: []string{arg_addr_str},
		UseAccess:       true,
	}
	err = atc.AddMethodCall(params)
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 1)
	txns, err := atc.BuildGroup()
	require.NoError(t, err)

	require.Empty(t, txns[0].Txn.ForeignApps)
	require.Empty(t, txns[0].Txn.ForeignAssets)
	require.Empty(t, txns[0].Txn.Accounts)

	require.Equal(t, len(txns[0].Txn.Access), 3)
	require.Equal(t, arg_addr, txns[0].Txn.Access[0].Address)
	require.Equal(t, types.AssetIndex(5), txns[0].Txn.Access[1].Asset)
	require.Equal(t, types.AppIndex(1), txns[0].Txn.Access[2].App)
}

func TestGatherSignatures(t *testing.T) {
	var atc AtomicTransactionComposer
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

	txWithSigners, _ := atc.BuildGroup()
	require.Equal(t, types.Digest{}, txWithSigners[0].Txn.Group)

	_, expectedSig, err := crypto.SignTransaction(account.PrivateKey, tx)
	require.NoError(t, err)
	require.Equal(t, len(sigs[0]), len(expectedSig))
	require.Equal(t, sigs[0], expectedSig)
}

func TestATCWithRejectVersion(t *testing.T) {
	var atc AtomicTransactionComposer
	account := crypto.GenerateAccount()

	method, err := abi.MethodFromSignature("add(uint64,uint64)uint64")
	require.NoError(t, err)

	sp := types.SuggestedParams{
		Fee:             1000,
		FirstRoundValid: 1000,
		LastRoundValid:  2000,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     []byte("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
		FlatFee:         true,
	}

	// Test with specific reject version
	err = atc.AddMethodCall(AddMethodCallParams{
		AppID:           123,
		Method:          method,
		MethodArgs:      []interface{}{uint64(1), uint64(2)},
		Sender:          account.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          BasicAccountTransactionSigner{Account: account},
		RejectVersion:   5,
	})
	require.NoError(t, err)

	txGroup, err := atc.BuildGroup()
	require.NoError(t, err)
	require.Equal(t, 1, len(txGroup))
	require.EqualValues(t, 5, txGroup[0].Txn.RejectVersion)

	// Test with default reject version (0)
	var atc2 AtomicTransactionComposer
	err = atc2.AddMethodCall(AddMethodCallParams{
		AppID:           456,
		Method:          method,
		MethodArgs:      []interface{}{uint64(3), uint64(4)},
		Sender:          account.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          BasicAccountTransactionSigner{Account: account},
	})
	require.NoError(t, err)

	txGroup2, err := atc2.BuildGroup()
	require.NoError(t, err)
	require.Equal(t, 1, len(txGroup2))
	require.EqualValues(t, 0, txGroup2[0].Txn.RejectVersion)

	// Test with UseAccess and reject version
	var atc3 AtomicTransactionComposer
	err = atc3.AddMethodCall(AddMethodCallParams{
		AppID:           789,
		Method:          method,
		MethodArgs:      []interface{}{uint64(5), uint64(6)},
		Sender:          account.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          BasicAccountTransactionSigner{Account: account},
		UseAccess:       true,
		RejectVersion:   7,
	})
	require.NoError(t, err)

	txGroup3, err := atc3.BuildGroup()
	require.NoError(t, err)
	require.Equal(t, 1, len(txGroup3))
	require.EqualValues(t, 7, txGroup3[0].Txn.RejectVersion)
}
