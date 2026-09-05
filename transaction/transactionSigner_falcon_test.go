//go:build falcon

package transaction

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
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
	require.Len(t, sigs, 1)

	var stx types.SignedTxn
	require.NoError(t, msgpack.Decode(sigs[0], &stx))
	require.Equal(t, tx, stx.Txn)
	require.Equal(t, types.Address{}, stx.AuthAddr)
	require.Equal(t, types.PQSchemeFalcon1024, stx.PQsig.Scheme)
	require.True(t, crypto.VerifyPQSig(transactionBytesToSign(tx), stx.PQsig))
}

func TestMakeFalcon1024AccountTransactionSignerWithAuthAddr(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)

	// Sender differs from the signer: the account has been rekeyed to the PQ key.
	fromAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	txSigner := PQAccountTransactionSigner{Signer: pqa.AsSigner()}
	tx := types.Transaction{
		Type:   types.PaymentTx,
		Header: types.Header{Sender: fromAddr, Fee: 217000, FirstValid: 972508, LastValid: 973508},
	}

	sigs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	var stx types.SignedTxn
	require.NoError(t, msgpack.Decode(sigs[0], &stx))
	require.Equal(t, pqa.Address(), stx.AuthAddr)
}

func TestMakeFalcon1024EmptyTransactionSigner(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	fromAddr := pqa.Address()
	toAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	txSigner := PQEmptyTransactionSigner{Signer: pqa.AsSigner()}
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     fromAddr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: toAddr,
			Amount:   5000,
		},
	}

	sigs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)
	require.Len(t, sigs, 1)

	var stx types.SignedTxn
	require.NoError(t, msgpack.Decode(sigs[0], &stx))
	require.Equal(t, tx, stx.Txn)
	require.Equal(t, types.Address{}, stx.AuthAddr)
	require.Equal(t, types.PQSchemeFalcon1024, stx.PQsig.Scheme)
	require.Equal(t, pqa.Salt, stx.PQsig.Salt)
	require.Equal(t, pqa.PublicKey[:], stx.PQsig.PublicKey)
	require.Empty(t, stx.PQsig.Signature)
	require.True(t, txSigner.Equals(PQEmptyTransactionSigner{Signer: pqa.AsSigner()}))
}

// The following golden tests are based on the PQ (Falcon-1024) fixtures from
// algorandfoundation/algokit-polytest. Each reconstructs the fixture
// transaction, re-signs it, and asserts the result equals the fixture's golden
// signed-transaction blob byte-for-byte.

// falcon1024GoldenAccount returns the fixed PQ account shared by all fixtures.
func falcon1024GoldenAccount(t *testing.T) crypto.Falcon1024Account {
	seed := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	pqa, err := crypto.Falcon1024AccountFromPQSeed(seed[:])
	require.NoError(t, err)
	return pqa
}

// falcon1024GoldenTxn builds the fixture payment transaction with the given sender.
func falcon1024GoldenTxn(t *testing.T, sender types.Address) types.Transaction {
	gh, err := base64.StdEncoding.DecodeString("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=")
	require.NoError(t, err)
	var genesisHash types.Digest
	copy(genesisHash[:], gh)
	// Payment to the all-zero address for 0 microAlgos, flat fee, no note/close.
	return types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:      sender,
			Fee:         1000,
			FirstValid:  50659540,
			LastValid:   50660540,
			GenesisID:   "testnet-v1.0",
			GenesisHash: genesisHash,
		},
	}
}

// Based on pqPayment.json.
func TestSignFalcon1024PaymentGolden(t *testing.T) {
	pqa := falcon1024GoldenAccount(t)

	// The derived address must match the fixture sender.
	require.Equal(t, "AZM6UV2ONIVHH7BK2CSBUPJCXNPZH5LFA2YFBCZPHSYXUFJ4LLLFJOUT5Y", pqa.Address().String())

	txn := falcon1024GoldenTxn(t, pqa.Address())
	_, stxBytes, err := SignTransaction(PQAccountTransactionSigner{Signer: pqa.AsSigner()}, txn)
	require.NoError(t, err)

	golden := "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"
	require.Equal(t, golden, base64.StdEncoding.EncodeToString(stxBytes))
}

// Based on pqRekeyedPayment.json.
func TestSignFalcon1024RekeyedPaymentGolden(t *testing.T) {
	pqa := falcon1024GoldenAccount(t)

	sender, err := types.DecodeAddress("BO6DI2SXMZ6DQAJAXWOH7V7FDUWF7X7KG7GS6W7UAWZMNP3PFV4J7HWYYY")
	require.NoError(t, err)

	txn := falcon1024GoldenTxn(t, sender)
	_, stxBytes, err := SignTransaction(PQAccountTransactionSigner{Signer: pqa.AsSigner()}, txn)
	require.NoError(t, err)

	golden := "g6VwcXNpZ4SicGvFBwEKg/0JiVcOB2K89BDC7T1drcUWJMSnZg35kZ4LGRoIWcM6YPAklosoCAufmazpmOMjvV1tuksbTd+lNrPL48EED4xOO88fCBNSQKvar0N44VBgOJVnF17AaFrTeDIE+2eJg7Gl+lYS5s0MvcyEUUbAxiZUU4Hr7eaWfi92OUKbGCVhtNIAwBg8JfnDV/c0kDFhCNHzEi2gMMlrJEjJnPKy3wVSBB73Ld7AcOhtqgOgCg0Q6ktdnN8NLDf39FGJQugUgFIZhRgAp/rXH1kWJqjFLq2kBggNs6YS/610IxtmU8wA2VJQzIGdG8Z4ZBsuQ5H0Oynv3kcxIdmDNg+rtPxos0uiqRIeDpq63iuSOti6YgGJpHH+1qFZQIlFmVcvScE9HHTTZ1JAokE64hRG71MmGbp7H0JpmCNLTW/4EEy/cMmIBY2K4/Ca0LqQeP11e5GWRzrMV05n7pdg0+vaCQdYKl33f05Vf+X7J7bFknJVmsC3X/JJ4k5DMl59huD0AybCbnA2e9UWqiFh4cqk2uVdlHHn20gM4QqyZOdYbyiOwFNnpEyvZCQwh5vVTS6KnjS8UvgkjQiZIV6IeSqIY8eBHolOY2R9lq5qStWEgN3YV8C4+cZaL2cg4HRzCJBSXeOZF2fBhOBD5+mEaBQZyfEK0G1IpUAFJ6K6m54N4q6iFUw9J4HOYHQ/UaJn6h6/d8U0R4aJ9R++1bJYEhNhdl9k9qrToIlImSiqCCMF/gYpFgEsyWn3blfdo+H9WL2NXKdKMdUNanNGMkBETybaExI1snQUUgqEQ9SOw+lnBX0bJOx4qOIGBxYXQL2CDkQraPs6BtQXgc0lj5l70dLpc1FQ4ZFii/ZRvtbmnsdOmmVzJPssMrwdMwEYk7i7SqxBGeCWD+1ZZF4QnSU022flWqJFAEUjZ1S0KrRoV8ooMzAc8kttWYAAcTbgwYTutM36o+XgH6AmqIqFW2oCIuezjgNMKv5os4bWF5plkgmLnBSRKDo3NAPQO6xQ4Q4HdIphkOUMpRigDdywWVOtiAIjEqTNYSlp6iVXDWqTWmxnLW9VLDb5ugFYzVJ1ScFRaqFEvDXlHduNqKnSUXiPhR2SoAW3jGO57h+r+RjMgh6+g2sdGrQ02PBVykMfWNCxi9enpGZjQm/n32mKlfSJo6M+pAGYN9yslsSIlpbauibGYoXYr+SAQi6zdxB3rRaZVqspU7R98RiFisWlsWF6odapR5kPZO8fQ0Q67lA1rytAQZIq8oLiq8VD2M+QGlJNVRiRAiWqYXQHao6bd+p6WsUi803NGyfY4oOKflMjVc65njjbN1JN94PavGlWocmqsG+4p59vn6U4dECdjLz0ARPkKdcGZC1XAEtyEp18J6ef0n7ZNgBHrP2+zlA0FxmWo7RjWe8OJo2YP2NF4BMzAIxyM2rFOExSUv5RtpJZEIEWH2VXtFUi2wLhEPK2q9JPrqsRwDDShmJzlmLZNq77U/J8cq+lhSkqyoFa7qdu1UpZdOg/WKZK4FstIhk6pZPz0OG85cqULMwkgbogpIxL3bStRzbNTH1Tosi77FzILX7eGkfYcQHBtTbXBKpQwUBWHC8wS05dIJTrCIKeIlvFhSCM9JV7DNtuSIrSis/LmzBIRvOSQeURCXKDHe9ekRqvN69K8CAUNLFQWhS8EicQSmXmleBZI0gSpwba/uZA031XnkWmGNG4H6RYFwtc6l9ypaklZ0hPriI2MUIRUdEDw95rdkjSvnNZhg8hOWvQqEThfkg3db2EBWp5V1cei1Wi6WXwPCxrI0JGktumIjpgr6pHKBIMScG1RdQgNAjRe2StoJNqWUOW7HYP/BCL29o9VLzqCzLeTJ8tyKreb+EHPqJQ9pC8XdQaCG7DzhvaevLWDrGemkkEjcH1EFCzwLqWiX79yAhCYti+CXS5H5Lgrum/+lnKvKIyuqidzvIJ2VPkh2iaF6zIZ24Mw2+WLFHnTIEW2uUALuD+RQi64tFC+F3YmmUWmqCPWT8HR94wTfm0CFROTVXGC2gtRpVAbh71W6CknLKuB0DC/Aoi/sY/F7g/gRLiyCsCsJIZEJXiBbn+ZGpVW52m0AWDT0GWqsFnQtsfNSgVHVZQKT4Y5fQGMWvldsGGhQpX277QC0oW3iijsQ90DBdGi9y89NaVRSXoy2hwWomujtLkLyiyLM5t4kec4bnheMpTm94i/SYxIxFgmrUJqpptBSa1I/RLDjFk7nBU5IGsDAW4mLPrVQLxEfHI4MaNjLeZtm51sdiURhcW4VpteEIIoUCNiluqAa4gdVzQxyXvADfoEiZHIHaY0YdDM02RkmJOtvhhABVKJaQ1zimoWboT72pNr5kTWZyIgK4keSdtBOgwIMAzZaxtdbKoM80Tk6NzY2jEAmYxo3NpZ8UE1roAV6UynngV1ymu1sdKy7HHidIVIwLMQcVyyQ9WOh0ydNREkL3LQttkk5krGY02TgY7QfhKSGmL0LWYwhklXOyVNHRoYUnL0naGbqGhULwMsNWbfFuwrrq++Jm4LanSz4Vfed3y2WhuV+XBtkYRaJQj7b207MwA1nahpK1aTnnHx/cS9y6NX6ZcRWUjWHysyGIxJByK1O41Fdb26r72KT5vd2bxH7vLJ/KkxIFgkEfYv/rK852YdzWClFcPLynfKonEHWOqvmixGaBfDwwq2zVpWkqOd3PpSNfzleny/f/Qp7UIYYpUa/P2UfE/GLUegKRar5Idz/btlNIi2J4WldOaeE45C4OiaR/ib1V7sE59fQVxKLfP/S3HOPjDEuKgmbE9zdOQeGT5xGhgM1NLfLYw/SJels3jUMWNc/FM9LIetEMAo2hGaM010ylkdrfGqLflYoMx0Qt8fjT3U070Vm9DVh7c2riF2NZXANqxhTJW1qGJTkjCr8rU3ZE5ss0W95BLGRQ3eE2Pnk8G7OFJocRmbhoqi+j1WBtU/TuYEQ5rJdPWL5a+EkEfaXi8rwRp1NBWsNOXeljFMWUpoTIH1hHvQs0kEfOgyqWnfZQ6iUsOQehusheLxGi6IiL0WPKeN4/B8eWtGBu85cQ2cWRBMd4znqknK7UXix+oDF99XMvCJvBI47DIPGfzNQ/CWFgagNC/2UHMcnQ6a371lhEUeWgoi2auUownxADVYtLEMIC4/WdzXkGeXkQNYmSm0rbhAYKr8HMY+FjcjxWLX9wgED1S49TjfZwLkzU2XRHRo2UUGZWf2YeSHM/7edlFsF0sGepZ65MLOzmQSHLPrMU17Z0IwbP2L17kMiBNsygzmpDSkOc7E6WXEf0z9oeiwxtCWxF30WHMsRVYCGqSA8VAnu/zriM8lRzlVhNC2SxUL1qu/qHmRctqGBShc6ItMjhlc4NDoOJN02bnVi1XnPt3ksb6pF7ZNhVU+NNpPYif1UnIXj4bMpydy/fGtIi62gLVm5/SGifLbT6zcCo6aDNeXnoG3ipK0K7O9k6LRxMYKRVQYj71p2GX5KpP3hTKGQypnjkpi6VcnlHlnQSBAOZt0IT+V7mDmWi/fWFfOkh/eyK9IoqjDlSl0GzR1qgmu6+sFdZkmO1PIQhX3oaWJYF3qDOiBjTODBWpw1HgBflhq5jtCihoUZlpxmn89m+f+qsIyjNdZ+bVnac4iw5FkIVDkeOvL1JOff0Tvvy1Gm/Cxa1UniJjCoQpSAGc6aFm0MvRY9zmFlnq78VkzX/eSxSvqEzNjYkYiGXdi2LK+8xpCX1P2xr0bODtOzquIc6ql9Ji0Db/AzTXE/70cE72aIxugNe3eBmLbb1KZRBDJJG4SDrMvsbUW6YKIuTrkljhlkjjiLoD703gZQ1LQcp00jLgb75esxCrL2/KMrIaYjbh6xqK3XUKMjQzm3oLlBTWGFWrdOQ2lb3KGtiOkyyHqh5G2LjdbE8Lx9rlpO1Siz2FJomSRsNSGLTPwNArkroxOLllEnsSAeN0DTx7EXTenkZGxrF4tEKk1HsseMb2Y8J27Bp3KZ6pTFSNVZ91Y09jt9wLMMVLqiS6ZTVJFlax/5UsMnkx0+jSiWE+t/Uyb0aLRDI0+Gotxipoo3NsdAOkc2ducsQgBlnqV05qKnP8KtCkGj0iu1+T9WUGsFCLLzyxehU8WtajdHhuh6NmZWXNA+iiZnbOAwUA1KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4DBQS8o3NuZMQgC7w0aldmfDgBIL2cf9flHSxf3+o3zS9b9AWyxr9vLXikdHlwZaNwYXk="
	require.Equal(t, golden, base64.StdEncoding.EncodeToString(stxBytes))
}
