//go:build falcon

package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func makeTestFalcon1024Account(t *testing.T) crypto.Falcon1024Account {
	mn := "auction inquiry lava second expand liberty glass involve ginger illness length room item discover ahead table doctor term tackle cement bonus profit right above catch"
	seed, err := mnemonic.ToPQSeed(mn, types.PQSchemeFalcon1024)
	require.NoError(t, err)
	pqa, err := crypto.Falcon1024AccountFromPQSeed(seed)
	require.NoError(t, err)
	return pqa
}

func TestMakeFalcon1024AccountTransactionSigner(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	fromAddr := pqa.Address()
	toAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	txSigner := PQAccountTransactionSigner{Signer: pqa.AsSigner()}
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

	sigs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	_, expectedSig, err := crypto.SignPQAccountTransaction(pqa.AsSigner(), tx)
	require.NoError(t, err)
	require.Equal(t, sigs[0], expectedSig)
}
