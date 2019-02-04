package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/crypto"
)

func TestEncodeDecode(t *testing.T) {
	a := Address{}
	for i := 0; i < 1000; i++ {
		crypto.RandomBytes(a[:])
		addr := a.String()
		b, err := DecodeAddress(addr)
		require.NoError(t, err)
		require.Equal(t, a, b)
	}
}

func TestGoldenValues(t *testing.T) {
	golden := "7777777777777777777777777777777777777777777777777774MSJUVU"
	a := Address{}
	for i := 0; i < len(a); i++ {
		a[i] = byte(0xFF)
	}
	require.Equal(t, golden, a.String())
}
