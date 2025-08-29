package crypto

import (
	"bytes"
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

func TestMultisigAccount_Address(t *testing.T) {
	addr1, err := types.DecodeAddress("XMHLMNAVJIMAW2RHJXLXKKK4G3J3U6VONNO3BTAQYVDC3MHTGDP3J5OCRU")
	require.NoError(t, err)
	addr2, err := types.DecodeAddress("HTNOX33OCQI2JCOLZ2IRM3BC2WZ6JUILSLEORBPFI6W7GU5Q4ZW6LINHLA")
	require.NoError(t, err)
	addr3, err := types.DecodeAddress("E6JSNTY4PVCY3IRZ6XEDHEO6VIHCQ5KGXCIQKFQCMB2N6HXRY4IB43VSHI")
	require.NoError(t, err)
	ma, err := MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})
	require.NoError(t, err)
	addrMultisig, err := ma.Address()
	require.NoError(t, err)
	require.Equal(t,
		"UCE2U2JC4O4ZR6W763GUQCG57HQCDZEUJY4J5I6VYY4HQZUJDF7AKZO5GM",
		addrMultisig.String(),
	)
}

func TestMultisigAccount_ZeroThreshInvalid(t *testing.T) {
	addr1, err := types.DecodeAddress("XMHLMNAVJIMAW2RHJXLXKKK4G3J3U6VONNO3BTAQYVDC3MHTGDP3J5OCRU")
	require.NoError(t, err)
	ma, err := MultisigAccountWithParams(1, 0, []types.Address{
		addr1,
	})
	require.Error(t, ma.Validate())
}

func TestMultisigAccount_Version1Only(t *testing.T) {
	addr1, err := types.DecodeAddress("XMHLMNAVJIMAW2RHJXLXKKK4G3J3U6VONNO3BTAQYVDC3MHTGDP3J5OCRU")
	require.NoError(t, err)
	ma, err := MultisigAccountWithParams(0, 1, []types.Address{
		addr1,
	})
	require.Error(t, ma.Validate())

	ma, err = MultisigAccountWithParams(2, 1, []types.Address{
		addr1,
	})
	require.Error(t, ma.Validate())
}

func TestLogicSigAddress(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	var args [][]byte

	expectedAddr, err := types.DecodeAddress("6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY")
	require.NoError(t, err)

	t.Run("no sig", func(t *testing.T) {
		var sk ed25519.PrivateKey
		var ma MultisigAccount

		lsig, err := makeLogicSig(program, args, sk, ma)
		require.NoError(t, err)

		actualAddr := LogicSigAddress(lsig)
		require.Equal(t, expectedAddr, actualAddr)
	})

	t.Run("single sig", func(t *testing.T) {
		account, err := AccountFromPrivateKey(ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42})
		require.NoError(t, err)

		var ma MultisigAccount

		lsig, err := makeLogicSig(program, args, account.PrivateKey, ma)
		require.NoError(t, err)

		// for backwards compatibility, we still expect the hashed program bytes address
		actualAddr := LogicSigAddress(lsig)
		require.Equal(t, expectedAddr, actualAddr)
	})

	t.Run("multi sig", func(t *testing.T) {
		ma, sk1, _, _ := makeTestMultisigAccount(t)

		lsig, err := makeLogicSig(program, args, sk1, ma)
		require.NoError(t, err)

		// for backwards compatibility, we still expect the hashed program bytes address
		actualAddr := LogicSigAddress(lsig)
		require.Equal(t, expectedAddr, actualAddr)
	})
}

func TestMakeLogicSigAccount(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	t.Run("Escrow", func(t *testing.T) {
		lsigAccount, err := MakeLogicSigAccountEscrowChecked(program, args)
		require.NoError(t, err)

		require.Equal(t, program, lsigAccount.Lsig.Logic)
		require.Equal(t, args, lsigAccount.Lsig.Args)
		require.Equal(t, types.Signature{}, lsigAccount.Lsig.Sig)
		require.True(t, lsigAccount.Lsig.Msig.Blank())
		require.Equal(t, ed25519.PublicKey(nil), lsigAccount.SigningKey)

		require.False(t, lsigAccount.IsDelegated())
	})

	t.Run("Delegated", func(t *testing.T) {
		account, err := AccountFromPrivateKey(ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42})
		require.NoError(t, err)

		lsigAccount, err := MakeLogicSigAccountDelegated(program, args, account.PrivateKey)
		require.NoError(t, err)

		expectedSignature := types.Signature{0x3e, 0x5, 0x3d, 0x39, 0x4d, 0xfb, 0x12, 0xbc, 0x65, 0x79, 0x9f, 0xea, 0x31, 0x8a, 0x7b, 0x8e, 0xa2, 0x51, 0x8b, 0x55, 0x2c, 0x8a, 0xbe, 0x6c, 0xd7, 0xa7, 0x65, 0x2d, 0xd8, 0xb0, 0x18, 0x7e, 0x21, 0x5, 0x2d, 0xb9, 0x24, 0x62, 0x89, 0x16, 0xe5, 0x61, 0x74, 0xcd, 0xf, 0x19, 0xac, 0xb9, 0x6c, 0x45, 0xa4, 0x29, 0x91, 0x99, 0x11, 0x1d, 0xe4, 0x7c, 0xe4, 0xfc, 0x12, 0xec, 0xce, 0x2}

		require.Equal(t, program, lsigAccount.Lsig.Logic)
		require.Equal(t, args, lsigAccount.Lsig.Args)
		require.Equal(t, expectedSignature, lsigAccount.Lsig.Sig)
		require.True(t, lsigAccount.Lsig.Msig.Blank())
		require.Equal(t, account.PublicKey, lsigAccount.SigningKey)

		require.True(t, lsigAccount.IsDelegated())
	})

	t.Run("DelegatedMsig", func(t *testing.T) {
		ma, sk1, sk2, _ := makeTestMultisigAccount(t)

		lsigAccount, err := MakeLogicSigAccountDelegatedMsig(program, args, ma, sk1)
		require.NoError(t, err)

		addr, err := ma.Address()
		require.NoError(t, err)

		var expectedSig types.Signature
		copy(expectedSig[:], ed25519.Sign(sk1, bytes.Join(
			[][]byte{[]byte("MsigProgram"), addr[:], program}, nil,
		)))

		expectedMsig := types.MultisigSig{
			Version:   ma.Version,
			Threshold: ma.Threshold,
			Subsigs: []types.MultisigSubsig{
				{
					Key: ma.Pks[0],
					Sig: expectedSig,
				},
				{
					Key: ma.Pks[1],
				},
				{
					Key: ma.Pks[2],
				},
			},
		}

		require.Equal(t, program, lsigAccount.Lsig.Logic)
		require.Equal(t, args, lsigAccount.Lsig.Args)
		require.Equal(t, types.Signature{}, lsigAccount.Lsig.Sig)
		require.Zero(t, lsigAccount.Lsig.Msig) // Legacy
		require.Equal(t, expectedMsig, lsigAccount.Lsig.LMsig)
		require.Nil(t, lsigAccount.SigningKey)

		require.True(t, lsigAccount.IsDelegated())

		t.Run("AppendMultisigSignature", func(t *testing.T) {
			err := lsigAccount.AppendMultisigSignature(sk2)
			require.NoError(t, err)

			var expectedSig types.Signature
			copy(expectedSig[:], ed25519.Sign(sk2, bytes.Join(
				[][]byte{[]byte("MsigProgram"), addr[:], program}, nil,
			)))
			expectedMsig.Subsigs[1].Sig = expectedSig

			require.Equal(t, program, lsigAccount.Lsig.Logic)
			require.Equal(t, args, lsigAccount.Lsig.Args)
			require.Equal(t, types.Signature{}, lsigAccount.Lsig.Sig)
			require.Zero(t, lsigAccount.Lsig.Msig) // Legacy
			require.Equal(t, expectedMsig, lsigAccount.Lsig.LMsig)
			require.Nil(t, lsigAccount.SigningKey)

			require.True(t, lsigAccount.IsDelegated())
		})
	})
}

func TestLogicSigAccountFromLogicSig(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	programAddr, err := types.DecodeAddress("6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY")
	require.NoError(t, err)

	programPublicKey := make(ed25519.PublicKey, len(programAddr))
	copy(programPublicKey, programAddr[:])

	t.Run("no sig", func(t *testing.T) {
		var sk ed25519.PrivateKey
		var ma MultisigAccount

		lsig, err := makeLogicSig(program, args, sk, ma)
		require.NoError(t, err)

		t.Run("with public key", func(t *testing.T) {
			_, err := LogicSigAccountFromLogicSig(lsig, &programPublicKey)
			require.Error(t, err, errLsigAccountPublicKeyNotNeeded)
		})

		t.Run("without public key", func(t *testing.T) {
			lsigAccount, err := LogicSigAccountFromLogicSig(lsig, nil)
			require.NoError(t, err)

			require.Equal(t, lsig, lsigAccount.Lsig)
			require.Equal(t, ed25519.PublicKey(nil), lsigAccount.SigningKey)

			require.False(t, lsigAccount.IsDelegated())
		})
	})

	t.Run("single sig", func(t *testing.T) {
		account, err := AccountFromPrivateKey(ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42})
		require.NoError(t, err)

		var ma MultisigAccount

		lsig, err := makeLogicSig(program, args, account.PrivateKey, ma)
		require.NoError(t, err)

		t.Run("with correct public key", func(t *testing.T) {
			lsigAccount, err := LogicSigAccountFromLogicSig(lsig, &account.PublicKey)
			require.NoError(t, err)

			require.Equal(t, lsig, lsigAccount.Lsig)
			require.Equal(t, account.PublicKey, lsigAccount.SigningKey)

			require.True(t, lsigAccount.IsDelegated())
		})

		t.Run("with incorrect public key", func(t *testing.T) {
			var wrongPublicKey = make(ed25519.PublicKey, len(account.PublicKey))
			copy(wrongPublicKey, account.PublicKey)
			wrongPublicKey[0] = 0xff
			_, err := LogicSigAccountFromLogicSig(lsig, &wrongPublicKey)
			require.Error(t, err, errLsigInvalidPublicKey)
		})

		t.Run("without public key", func(t *testing.T) {
			_, err := LogicSigAccountFromLogicSig(lsig, nil)
			require.Error(t, err, errLsigNoPublicKey)
		})
	})

	t.Run("multi sig", func(t *testing.T) {
		ma, sk1, _, _ := makeTestMultisigAccount(t)

		lsig, err := makeLogicSig(program, args, sk1, ma)
		require.NoError(t, err)

		t.Run("with public key", func(t *testing.T) {
			msigAddr, err := ma.Address()
			require.NoError(t, err)

			msigPublicKey := make(ed25519.PublicKey, len(msigAddr))
			copy(msigPublicKey, msigAddr[:])

			_, err = LogicSigAccountFromLogicSig(lsig, &msigPublicKey)
			require.Error(t, err, errLsigAccountPublicKeyNotNeeded)
		})

		t.Run("without public key", func(t *testing.T) {
			lsigAccount, err := LogicSigAccountFromLogicSig(lsig, nil)
			require.NoError(t, err)

			require.Equal(t, lsig, lsigAccount.Lsig)
			require.Equal(t, ed25519.PublicKey(nil), lsigAccount.SigningKey)

			require.True(t, lsigAccount.IsDelegated())
		})
	})
}

func TestLogicSigAccount_Address(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	t.Run("no sig", func(t *testing.T) {
		lsigAccount, err := MakeLogicSigAccountEscrowChecked(program, args)
		require.NoError(t, err)

		expectedAddr, err := types.DecodeAddress("6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY")
		require.NoError(t, err)

		actualAddr, err := lsigAccount.Address()
		require.NoError(t, err)
		require.Equal(t, expectedAddr, actualAddr)
	})

	t.Run("single sig", func(t *testing.T) {
		account, err := AccountFromPrivateKey(ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42})
		require.NoError(t, err)

		lsigAccount, err := MakeLogicSigAccountDelegated(program, args, account.PrivateKey)
		require.NoError(t, err)

		actualAddr, err := lsigAccount.Address()
		require.NoError(t, err)
		require.Equal(t, account.Address, actualAddr)
	})

	t.Run("multi sig", func(t *testing.T) {
		ma, sk1, _, _ := makeTestMultisigAccount(t)
		maAddr, err := ma.Address()
		require.NoError(t, err)

		lsigAccount, err := MakeLogicSigAccountDelegatedMsig(program, args, ma, sk1)
		require.NoError(t, err)

		actualAddr, err := lsigAccount.Address()
		require.NoError(t, err)
		require.Equal(t, maAddr, actualAddr)
	})
}
