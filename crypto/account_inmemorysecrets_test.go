package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func TestGenerateAccount(t *testing.T) {
	kp := GenerateAccount()

	// Public key should not be empty
	require.NotEqual(t, ed25519.PublicKey{}, kp.PublicKey)

	// Private key should not be empty
	require.NotEqual(t, ed25519.PrivateKey{}, kp.PrivateKey)

	// Account should equal itself
	require.Equal(t, kp, kp)

	// Address should be identical to public key
	pk := ed25519.PublicKey(kp.Address[:])
	require.Equal(t, pk, kp.PublicKey)

	message := []byte("test message")
	sig := ed25519.Sign(kp.PrivateKey, message)
	// Public key should verify signature from private key
	require.True(t, ed25519.Verify(kp.PublicKey, message, sig))

	kp2 := GenerateAccount()
	// Calling the function again should produce a different account
	require.NotEqual(t, kp, kp2)
}

func TestAccountFromPrivateKey(t *testing.T) {
	exampleAccount := Account{
		PrivateKey: ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42},
		PublicKey:  ed25519.PublicKey{0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42},
		Address:    types.Address{0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42},
	}

	t.Run("From private key", func(t *testing.T) {
		pk := exampleAccount.PrivateKey[:]

		actual, err := AccountFromPrivateKey(pk)
		require.NoError(t, err)

		require.Equal(t, exampleAccount, actual)
	})

	t.Run("From seed only", func(t *testing.T) {
		pk := exampleAccount.PrivateKey.Seed() // get just the seed portion of the private key (first 32 bytes)

		_, err := AccountFromPrivateKey(pk)
		require.Error(t, err, errInvalidPrivateKey)
	})

	t.Run("From mnemonic", func(t *testing.T) {
		m := "olympic cricket tower model share zone grid twist sponsor avoid eight apology patient party success claim famous rapid donor pledge bomb mystery security ability often"
		pk, err := mnemonic.ToPrivateKey(m)
		require.NoError(t, err)

		actual, err := AccountFromPrivateKey(pk)
		require.NoError(t, err)

		require.Equal(t, exampleAccount, actual)
	})
}

func TestAccountAsSigner(t *testing.T) {
	acct := GenerateAccount()

	sgnr := acct.AsSigner()
	require.Equal(t, Ed25519PublicKey(acct.Address), sgnr.Ed25519PublicKey())

	message := []byte("test message")
	sig, err := sgnr.Ed25519Sign(message)
	require.NoError(t, err)
	require.True(t, ed25519.Verify(acct.PublicKey, message, sig))
}

// The deprecated in-memory API is a thin wrapper around the Ed25519 API, which
// is where the behaviour itself is tested. These tests only check that every
// wrapper delegates to its replacement.

func TestMakeLogicSigAccountDelegatedInMemory(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{{0x01}, {0x02, 0x03}}
	acct := GenerateAccount()

	lsa, err := MakeLogicSigAccountDelegated(program, args, acct.PrivateKey)
	require.NoError(t, err)

	expectedLsa, err := Ed25519MakeLogicSigAccountDelegated(program, args, acct.AsSigner())
	require.NoError(t, err)
	require.Equal(t, expectedLsa, lsa)
}

func TestMakeLogicSigAccountDelegatedMsigInMemory(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{{0x01}, {0x02, 0x03}}
	ma, acct1, acct2 := makeInMemoryTestMultisigAccount(t)

	lsa, err := MakeLogicSigAccountDelegatedMsig(program, args, ma, acct1.PrivateKey)
	require.NoError(t, err)

	expectedLsa, err := Ed25519MakeLogicSigAccountDelegatedMsig(program, args, ma, acct1.AsSigner())
	require.NoError(t, err)
	require.Equal(t, expectedLsa, lsa)

	// AppendMultisigSignature completes the 2-of-2 delegation
	require.NoError(t, lsa.AppendMultisigSignature(acct2.PrivateKey))
	require.NoError(t, expectedLsa.Ed25519AppendMultisigSignature(acct2.AsSigner()))
	require.Equal(t, expectedLsa, lsa)

	maAddr, err := ma.Address()
	require.NoError(t, err)
	require.True(t, VerifyLogicSig(lsa.Lsig, maAddr))

	addr, err := lsa.Address()
	require.NoError(t, err)
	require.Equal(t, maAddr, addr)
}
