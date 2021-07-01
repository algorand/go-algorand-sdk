package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

func TestKeyGeneration(t *testing.T) {
	kp := GenerateAccount()

	// Public key should not be empty
	require.NotEqual(t, kp.PublicKey, ed25519.PublicKey{})

	// Public key should not be empty
	require.NotEqual(t, kp.PrivateKey, ed25519.PrivateKey{})

	// Address should be identical to public key
	pk := ed25519.PublicKey(kp.Address[:])
	require.Equal(t, pk, kp.PublicKey)
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
