package mnemonic

import (
	"crypto/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func TestGenerateAndRecovery(t *testing.T) {
	key := make([]byte, 32)
	for i := 0; i < 1000; i++ {
		// Generate a key
		_, err := rand.Read(key)
		require.NoError(t, err)
		// Go from key -> mnemonic
		m, err := FromKey(key)
		// Go from mnemonic -> key
		recovered, err := ToKey(m)
		require.NoError(t, err)
		require.Equal(t, recovered, key)
	}

	var mdk types.MasterDerivationKey
	for i := 0; i < 1000; i++ {
		// Generate a mdk
		_, err := rand.Read(key)
		require.NoError(t, err)
		copy(mdk[:], key)
		// Go from mdk -> mnemonic
		m, err := FromMasterDerivationKey(mdk)
		// Go from mnemonic -> mdk
		recovered, err := ToMasterDerivationKey(m)
		require.NoError(t, err)
		require.Equal(t, recovered, mdk)
	}

	for i := 0; i < 1000; i++ {
		// Generate a private key
		acct := crypto.GenerateAccount()
		// Go from sk -> mnemonic
		m, err := FromPrivateKey(acct.PrivateKey)
		// Go from mnemonic -> sk
		recovered, err := ToPrivateKey(m)
		require.NoError(t, err)
		require.Equal(t, recovered, acct.PrivateKey)
	}
}

func TestZeroVector(t *testing.T) {
	zeroVector := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	mn := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon invest"

	m, err := FromKey(zeroVector)
	require.NoError(t, err)
	require.Equal(t, mn, m)
	return
}

func TestWordNotInList(t *testing.T) {
	mn := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon zzz invest"
	_, err := ToKey(mn)
	require.Error(t, err)
	return
}

func TestCorruptedChecksum(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	m, err := FromKey(key)
	wl := strings.Split(m, sepStr)
	lastWord := wl[len(wl)-1]
	// Shuffle the last word (last 11 bits of checksum)
	wl[len(wl)-1] = wordlist[(indexOf(wordlist, lastWord)+1)%len(wordlist)]
	recovered, err := ToKey(strings.Join(wl, sepStr))
	require.Error(t, err)
	require.Empty(t, recovered)
}

func TestInvalidKeyLen(t *testing.T) {
	badLens := []int{0, 31, 33, 100}
	for _, l := range badLens {
		key := make([]byte, l)
		_, err := rand.Read(key)
		require.NoError(t, err)
		m, err := FromKey(key)
		require.Error(t, err)
		require.Empty(t, m)
	}
}
