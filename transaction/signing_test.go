package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func makeSigningTestTxn(sender, receiver types.Address) types.Transaction {
	return types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     sender,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: receiver,
			Amount:   5000,
		},
	}
}

// The signers in this package assemble signed transactions out of the crypto
// package's primitives, while crypto keeps its own deprecated in-memory
// entrypoints. Pin the two to the same bytes, so that neither can drift into
// signing something the other does not.
func TestSigningMatchesDeprecatedCryptoHelpers(t *testing.T) {
	account := crypto.GenerateAccount()
	other := crypto.GenerateAccount()

	t.Run("ed25519", func(t *testing.T) {
		for name, sender := range map[string]types.Address{
			"sender is signer": account.Address,
			"rekeyed":          other.Address,
		} {
			t.Run(name, func(t *testing.T) {
				tx := makeSigningTestTxn(sender, other.Address)

				expectedTxid, expectedStx, err := crypto.SignTransaction(account.PrivateKey, tx) //nolint:staticcheck // pinning the deprecated path
				require.NoError(t, err)

				txid, stx, err := SignTransaction(Ed25519AccountTransactionSigner{Signer: account.AsSigner()}, tx)
				require.NoError(t, err)

				require.Equal(t, expectedTxid, txid)
				require.Equal(t, expectedStx, stx)
			})
		}
	})

	t.Run("multisig", func(t *testing.T) {
		ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{account.Address, other.Address})
		require.NoError(t, err)
		maAddr, err := ma.Address()
		require.NoError(t, err)

		tx := makeSigningTestTxn(maAddr, other.Address)

		expectedTxid, expectedStx, err := crypto.SignMultisigTransaction(account.PrivateKey, ma, tx) //nolint:staticcheck // pinning the deprecated path
		require.NoError(t, err)

		txid, stx, err := SignTransaction(MultiSigEd25519AccountTransactionSigner{
			Msig:    ma,
			Signers: []crypto.Ed25519Signer{account.AsSigner()},
		}, tx)
		require.NoError(t, err)

		require.Equal(t, expectedTxid, txid)
		require.Equal(t, expectedStx, stx)

		// ...and the same for appending the second signature to that partial blob.
		expectedTxid, expectedStx, err = crypto.AppendMultisigTransaction(other.PrivateKey, ma, expectedStx) //nolint:staticcheck // pinning the deprecated path
		require.NoError(t, err)

		txid, stx, err = Ed25519AccountTransactionSigner{Signer: other.AsSigner()}.AppendSignature(ma, stx)
		require.NoError(t, err)

		require.Equal(t, expectedTxid, txid)
		require.Equal(t, expectedStx, stx)
	})
}

// The bytes a transaction signature commits to are defined by crypto; make sure
// this package signs exactly those.
func TestTransactionBytesToSignIsSharedWithCrypto(t *testing.T) {
	account := crypto.GenerateAccount()
	tx := makeSigningTestTxn(account.Address, account.Address)

	toBeSigned := crypto.TransactionBytesToSign(tx)
	require.Equal(t, []byte("TX"), toBeSigned[:2])

	expected, err := crypto.Ed25519RawSignature(account.AsSigner(), toBeSigned)
	require.NoError(t, err)

	stxBytes, err := ed25519SignTransaction(account.AsSigner(), tx)
	require.NoError(t, err)

	var stx types.SignedTxn
	require.NoError(t, msgpack.Decode(stxBytes, &stx))
	require.Equal(t, expected, stx.Sig)
}
