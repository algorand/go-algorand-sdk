package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"

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

// makeInMemoryTestMultisigAccount returns a 2-of-2 multisig account along with
// the in-memory accounts backing it
func makeInMemoryTestMultisigAccount(t *testing.T) (MultisigAccount, Account, Account) {
	acct1 := GenerateAccount()
	acct2 := GenerateAccount()
	ma, err := MultisigAccountWithParams(1, 2, []types.Address{acct1.Address, acct2.Address})
	require.NoError(t, err)
	return ma, acct1, acct2
}

func TestGenerateAddressFromSK(t *testing.T) {
	acct := GenerateAccount()

	addr, err := GenerateAddressFromSK(acct.PrivateKey)
	require.NoError(t, err)
	require.Equal(t, acct.Address, addr)
}

func TestSKToInMemorySigner(t *testing.T) {
	acct := GenerateAccount()

	sgnr, err := SKToInMemorySigner(acct.PrivateKey)
	require.NoError(t, err)
	require.Equal(t, Ed25519PublicKey(acct.Address), sgnr.Ed25519PublicKey())

	message := []byte("test message")
	sig, err := sgnr.Ed25519Sign(message)
	require.NoError(t, err)
	require.True(t, ed25519.Verify(acct.PublicKey, message, sig))
}

func TestInMemorySecretsDelegateToEd25519Signer(t *testing.T) {
	ma, acct1, acct2 := makeInMemoryTestMultisigAccount(t)
	sgnr1, err := SKToInMemorySigner(acct1.PrivateKey)
	require.NoError(t, err)
	sgnr2, err := SKToInMemorySigner(acct2.PrivateKey)
	require.NoError(t, err)

	program := []byte{1, 32, 1, 1, 34}
	data := []byte{0x01, 0x02, 0x03}

	t.Run("SignTransaction", func(t *testing.T) {
		tx := makeInMemoryTestTxn(acct1.Address)

		txid, stx, err := SignTransaction(acct1.PrivateKey, tx)
		require.NoError(t, err)

		expectedTxid, expectedStx, err := Ed25519SignTransaction(sgnr1, tx)
		require.NoError(t, err)
		require.Equal(t, expectedTxid, txid)
		require.Equal(t, expectedStx, stx)
	})

	t.Run("SignBytes", func(t *testing.T) {
		sig, err := SignBytes(acct1.PrivateKey, data)
		require.NoError(t, err)

		expectedSig, err := Ed25519SignBytes(sgnr1, data)
		require.NoError(t, err)
		require.Equal(t, expectedSig, sig)
	})

	t.Run("SignBid", func(t *testing.T) {
		bid := types.Bid{
			BidderKey:   acct1.Address,
			BidCurrency: 1000,
			MaxPrice:    10,
			BidID:       1,
			AuctionKey:  acct2.Address,
			AuctionID:   2,
		}

		signedBid, err := SignBid(acct1.PrivateKey, bid)
		require.NoError(t, err)

		expectedSignedBid, err := Ed25519SignBid(sgnr1, bid)
		require.NoError(t, err)
		require.Equal(t, expectedSignedBid, signedBid)
	})

	t.Run("SignMultisigTransaction", func(t *testing.T) {
		maAddr, err := ma.Address()
		require.NoError(t, err)
		tx := makeInMemoryTestTxn(maAddr)

		txid, stx, err := SignMultisigTransaction(acct1.PrivateKey, ma, tx)
		require.NoError(t, err)

		expectedTxid, expectedStx, err := Ed25519SignMultisigTransaction(sgnr1, ma, tx)
		require.NoError(t, err)
		require.Equal(t, expectedTxid, txid)
		require.Equal(t, expectedStx, stx)
	})

	t.Run("AppendMultisigTransaction", func(t *testing.T) {
		maAddr, err := ma.Address()
		require.NoError(t, err)
		tx := makeInMemoryTestTxn(maAddr)
		_, preStx, err := Ed25519SignMultisigTransaction(sgnr1, ma, tx)
		require.NoError(t, err)

		txid, stx, err := AppendMultisigTransaction(acct2.PrivateKey, ma, preStx)
		require.NoError(t, err)

		expectedTxid, expectedStx, err := Ed25519AppendMultisigTransaction(sgnr2, ma, preStx)
		require.NoError(t, err)
		require.Equal(t, expectedTxid, txid)
		require.Equal(t, expectedStx, stx)
	})

	t.Run("AppendMultisigToLogicSig", func(t *testing.T) {
		// ed25519 signatures are deterministic, so both LogicSigs start out identical
		makeLsig := func() types.LogicSig {
			lsa, err := Ed25519MakeLogicSigAccountDelegatedMsig(program, nil, ma, sgnr1)
			require.NoError(t, err)
			return lsa.Lsig
		}

		lsig := makeLsig()
		require.NoError(t, AppendMultisigToLogicSig(&lsig, acct2.PrivateKey))

		expectedLsig := makeLsig()
		require.NoError(t, Ed25519AppendMultisigToLogicSig(&expectedLsig, sgnr2))
		require.Equal(t, expectedLsig, lsig)
	})

	t.Run("TealSign", func(t *testing.T) {
		contractAddress := AddressFromProgram(program)

		sig, err := TealSign(acct1.PrivateKey, data, contractAddress)
		require.NoError(t, err)

		expectedSig, err := Ed25519TealSign(sgnr1, data, contractAddress)
		require.NoError(t, err)
		require.Equal(t, expectedSig, sig)
	})

	t.Run("TealSignFromProgram", func(t *testing.T) {
		sig, err := TealSignFromProgram(acct1.PrivateKey, data, program)
		require.NoError(t, err)

		expectedSig, err := Ed25519TealSignFromProgram(sgnr1, data, program)
		require.NoError(t, err)
		require.Equal(t, expectedSig, sig)
	})
}
