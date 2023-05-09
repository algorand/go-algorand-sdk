package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	json2 "github.com/algorand/go-algorand-sdk/v2/encoding/json"
)

func TestUnmarshalBlockHash(t *testing.T) {
	testcases := []struct {
		name      string
		input     string
		outputB64 string
		err       string
	}{
		{
			name:      "blk-B32",
			input:     "blk-PEMLJXJYPLIFJ7DGVGTSEGUHUFR3M6Y67UZW3AC4LLNLF26XIXQA",
			outputB64: "eRi03Th60FT8ZqmnIhqHoWO2ex79M22AXFrasuvXReA=",
		}, {
			name:      "B32",
			input:     "PEMLJXJYPLIFJ7DGVGTSEGUHUFR3M6Y67UZW3AC4LLNLF26XIXQA",
			outputB64: "eRi03Th60FT8ZqmnIhqHoWO2ex79M22AXFrasuvXReA=",
		}, {
			name:      "B64",
			input:     "eRi03Th60FT8ZqmnIhqHoWO2ex79M22AXFrasuvXReA=",
			outputB64: "eRi03Th60FT8ZqmnIhqHoWO2ex79M22AXFrasuvXReA=",
		}, {
			name:  "B64-err-length",
			input: "AAE=",
			err:   "decoded block hash is the wrong length",
		}, {
			name:  "B64-err-illegal",
			input: "bogus",
			err:   "illegal base64 data at input byte 4",
		}, {
			name:  "Overflow does not panic",
			input: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ",
			err:   "illegal base64 data at input byte 56",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// wrap in quotes so that it is valid JSON
			data := []byte(fmt.Sprintf("\"%s\"", tc.input))

			var hash BlockHash
			err := json2.Decode(data, &hash)
			if tc.err != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			actual, err := hash.MarshalText()
			require.NoError(t, err)
			require.Equal(t, tc.outputB64, string(actual))
		})
	}
}
