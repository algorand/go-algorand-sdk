package crypto

import (
	"github.com/algorand/go-algorand-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyGeneration(t *testing.T) {
	var sk, pk []byte
	addr := GenerateAccount(sk, pk)

	// Public key should not be empty
	require.NotEqual(t, pk, []byte{})

	// Public key should not be empty
	require.NotEqual(t, sk, []byte{})

	// Address should be identical to public key
	decoded, err := types.DecodeAddress(addr)
	if err != nil {
		return
	}
	require.Equal(t, pk, []byte(decoded[:]))
}
