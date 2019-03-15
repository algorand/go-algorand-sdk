package crypto

import (
	"github.com/algorand/go-algorand-sdk/types"
	"golang.org/x/crypto/ed25519"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyGenerationSplit(t *testing.T) {
	sk := GenerateSK()

	// Public key should not be empty
	require.NotEqual(t, sk, []byte{})

}

func TestKeyGenerationSAddress(t *testing.T) {
	sk := GenerateSK()

	// Public key should not be empty
	require.NotEqual(t, sk, []byte{})

	addr := GenerateAddressFromSK(sk)

	require.Len(t, addr, 58)

	// Address should be identical to public key
	decoded, err := types.DecodeAddress(addr)
	if err != nil {
		return
	}

	require.Equal(t, ed25519.PrivateKey(sk).Public().(ed25519.PublicKey), ed25519.PublicKey(decoded[:]))

}
