package future

import (
	"crypto/ed25519"
	"testing"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMakeBasicAccountTransactionSigner(t *testing.T) {
	account := crypto.GenerateAccount()
	txSigner := MakeBasicAccountTransactionSigner(account)

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

	sigs, _, err := txSigner([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	_, expectedSig, err := crypto.SignTransaction(account.PrivateKey, tx)
	require.NoError(t, err)
	require.Equal(t, sigs[0], expectedSig)
}

func TestMakeLogicSigAccountTransactionSigner(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}
	account := crypto.GenerateAccount()
	lsig, err := crypto.MakeLogicSigAccountDelegated(program, args, account.PrivateKey)
	require.NoError(t, err)

	programHash := "6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY"
	programAddr, err := types.DecodeAddress(programHash)
	require.NoError(t, err)

	txSigner := MakeLogicSigAccountTransactionSigner(lsig)

	require.NoError(t, err)
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     programAddr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: programAddr,
			Amount:   5000,
		},
	}

	sigs, _, err := txSigner([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	_, expectedSig, err := crypto.SignLogicSigAccountTransaction(lsig, tx)
	require.NoError(t, err)
	require.Equal(t, sigs[0], expectedSig)
}

func makeTestMultisigAccount(t *testing.T) (crypto.MultisigAccount, ed25519.PrivateKey, ed25519.PrivateKey, ed25519.PrivateKey) {
	addr1, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)
	addr2, err := types.DecodeAddress("BFRTECKTOOE7A5LHCF3TTEOH2A7BW46IYT2SX5VP6ANKEXHZYJY77SJTVM")
	require.NoError(t, err)
	addr3, err := types.DecodeAddress("47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU")
	require.NoError(t, err)
	ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})
	require.NoError(t, err)
	mn1 := "auction inquiry lava second expand liberty glass involve ginger illness length room item discover ahead table doctor term tackle cement bonus profit right above catch"
	sk1, err := mnemonic.ToPrivateKey(mn1)
	require.NoError(t, err)
	mn2 := "since during average anxiety protect cherry club long lawsuit loan expand embark forum theory winter park twenty ball kangaroo cram burst board host ability left"
	sk2, err := mnemonic.ToPrivateKey(mn2)
	require.NoError(t, err)
	mn3 := "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	sk3, err := mnemonic.ToPrivateKey(mn3)
	return ma, sk1, sk2, sk3
}

func TestMakeMultiSigAccountTransactionSigner(t *testing.T) {
	ma, sk1, _, _ := makeTestMultisigAccount(t)
	fromAddr, err := ma.Address()
	require.NoError(t, err)
	toAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	txSigner := MakeMultiSigAccountTransactionSigner(ma, [][]byte{sk1})
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     fromAddr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: toAddr,
			Amount:   5000,
		},
	}

	sigs, _, err := txSigner([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	_, expectedSig, err := crypto.SignMultisigTransaction(sk1, ma, tx)
	require.NoError(t, err)
	require.Equal(t, sigs[0], expectedSig)
}

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
	txSigner := MakeBasicAccountTransactionSigner(account)

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
	txSigner := MakeBasicAccountTransactionSigner(account)

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

	_, err = atc.BuildGroup()
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILT)

	err = atc.AddTransaction(txAndSigner)
	require.Error(t, err)
}

func TestAddTransactionWithMaxTransactions(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := MakeBasicAccountTransactionSigner(account)

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
	txSigner := MakeBasicAccountTransactionSigner(account)
	methodSig := "add()uint32"

	method, err := MethodFromSignature(methodSig)
	require.NoError(t, err)

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	err = atc.AddMethodCall(
		0,
		method,
		[]MethodArgument{},
		addr,
		types.SuggestedParams{},
		types.NoOpOC,
		[]byte{},
		[32]byte{},
		addr,
		txSigner,
	)
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), BUILDING)
	require.Equal(t, atc.Count(), 1)
}

func TestGatherSignatures(t *testing.T) {
	atc := MakeAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := MakeBasicAccountTransactionSigner(account)

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
	_, expectedSig, err := crypto.SignTransaction(account.PrivateKey, tx)
	require.NoError(t, err)
	require.Equal(t, sigs[0], expectedSig)
}
