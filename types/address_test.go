package types

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	json2 "github.com/algorand/go-algorand-sdk/v2/encoding/json"
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

		addr, err := EncodeAddress(a[:])
		require.NoError(t, err)

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

func TestUnmarshalAddress(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		str    string
		output string
		err    string
	}{
		{
			name:   "B32+Checksum",
			input:  "7HJBGRIWI7GDL42SOJNIAZ7LJ7EBEGKGE5S52QZXAWDXOHDKMDFR6AUXDE",
			str:    "7HJBGRIWI7GDL42SOJNIAZ7LJ7EBEGKGE5S52QZXAWDXOHDKMDFR6AUXDE",
			output: "+dITRRZHzDXzUnJagGfrT8gSGUYnZd1DNwWHdxxqYMs=",
		}, {
			name:   "B64",
			input:  "+dITRRZHzDXzUnJagGfrT8gSGUYnZd1DNwWHdxxqYMs=",
			str:    "7HJBGRIWI7GDL42SOJNIAZ7LJ7EBEGKGE5S52QZXAWDXOHDKMDFR6AUXDE",
			output: "+dITRRZHzDXzUnJagGfrT8gSGUYnZd1DNwWHdxxqYMs=",
		}, {
			name:  "B64-err-length",
			input: "AAE=",
			err:   "decoded address is the wrong length",
		}, {
			name:  "B64-err-illegal",
			input: "bogus",
			err:   "illegal base64 data at input byte 4",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// wrap in quotes so that it is valid JSON
			data := []byte(fmt.Sprintf("\"%s\"", tc.input))

			var addr Address
			err := json2.Decode(data, &addr)
			if tc.err != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.str, addr.String())
			actual, err := addr.MarshalText()
			require.NoError(t, err)
			require.Equal(t, tc.output, string(actual))
		})
	}
}

func TestDecodeNonCanonicalAddress(t *testing.T) {
	// Canonical addresses must end with one of the following: "AEIMQUY4",
	// e.g. "7HJBGRIWI7GDL42SOJNIAZ7LJ7EBEGKGE5S52QZXAWDXOHDKMDFR6AUXDE"
	addrs := []string{
		"7HJBGRIWI7GDL42SOJNIAZ7LJ7EBEGKGE5S52QZXAWDXOHDKMDFR6AUXDF",
		"7HJBGRIWI7GDL42SOJNIAZ7LJ7EBEGKGE5S52QZXAWDXOHDKMDFR6AUXDG",
		"7HJBGRIWI7GDL42SOJNIAZ7LJ7EBEGKGE5S52QZXAWDXOHDKMDFR6AUXDH",
	}
	for _, addr := range addrs {
		_, err := DecodeAddress(addr)
		require.Error(t, err)
		require.ErrorContains(t, err, fmt.Sprintf("address %s is non-canonical", addr))
	}
}
