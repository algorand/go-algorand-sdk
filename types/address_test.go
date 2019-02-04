package types

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomBytes(s []byte) {
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	a := Address{}
	for i := 0; i < 1000; i++ {
		randomBytes(a[:])
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
