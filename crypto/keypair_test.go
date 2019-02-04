package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

func TestKeyGeneration(t *testing.T) {
	kp := GenerateKeyPair()

	// Public key should not be empty
	require.NotEqual(t, kp.PublicKey, ed25519.PublicKey{})

	// Public key should not be empty
	require.NotEqual(t, kp.PrivateKey, ed25519.PrivateKey{})

	// Address should be identical to public key
	pk := ed25519.PublicKey(kp.Address[:])
	require.Equal(t, pk, kp.PublicKey)
}
