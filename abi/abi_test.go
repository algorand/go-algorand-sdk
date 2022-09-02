package abi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportWorks(t *testing.T) {
	// This test is not meant to be exhaustive. It's just a simple test to
	// verify importing the ABI package from avm-abi is working

	abiType, err := TypeOf("uint64")
	require.NoError(t, err)

	valueToEncode := 10000
	expected := []byte{0, 0, 0, 0, 0, 0, 39, 16}

	actual, err := abiType.Encode(valueToEncode)
	require.NoError(t, err)

	require.Equal(t, expected, actual)
}
