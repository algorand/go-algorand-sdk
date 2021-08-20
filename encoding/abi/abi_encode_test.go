package abi

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeValid(t *testing.T) {
	for intSize := 8; intSize <= 512; intSize += 8 {
		upperLimit := big.NewInt(0).Lsh(big.NewInt(1), uint(intSize))
		for i := 0; i < 1000; i++ {
			randomInt, err := rand.Int(rand.Reader, upperLimit)
			require.NoError(t, err, "cryptographic random int init fail")
			valueUint, err := MakeUint(randomInt, uint16(intSize))
			require.NoError(t, err, "makeUint Fail")
			encodedUint, err := valueUint.Encode()
			require.NoError(t, err, "uint encode fail")
			buffer := make([]byte, intSize/8)
			randomInt.FillBytes(buffer)
			require.Equal(t, buffer, encodedUint, "encode uint not match with expected")
		}
	}

	for size := 8; size <= 512; size += 8 {
		upperLimit := big.NewInt(0).Lsh(big.NewInt(1), uint(size))
		for precision := 1; precision <= 160; precision++ {
			denomLimit := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil)
			for i := 0; i < 10; i++ {
				randomInt, err := rand.Int(rand.Reader, upperLimit)
				require.NoError(t, err, "cryptographic random int init fail")

				ufixedRational := big.NewRat(1, 1).SetFrac(randomInt, denomLimit)
				valueUfixed, err := MakeUfixed(ufixedRational, uint16(size), uint16(precision))
				require.NoError(t, err, "makeUfixed Fail")

				encodedUfixed, err := valueUfixed.Encode()
				require.NoError(t, err, "ufixed encode fail")

				buffer := make([]byte, size/8)
				randomInt.FillBytes(buffer)
				require.Equal(t, buffer, encodedUfixed, "encode ufixed not match with expected")
			}
		}
	}

	upperLimit := big.NewInt(0).Lsh(big.NewInt(1), 256)
	for i := 0; i < 1000; i++ {
		randomAddrInt, err := rand.Int(rand.Reader, upperLimit)
		require.NoError(t, err, "cryptographic random int init fail")

		address := make([]byte, 32)
		randomAddrInt.FillBytes(address)

		var addrBytes [32]byte
		copy(addrBytes[:], address[:32])

		addressValue := MakeAddress(addrBytes)
		addrEncode, err := addressValue.Encode()
		require.NoError(t, err, "address encode fail")
		require.Equal(t, address, addrEncode, "encode addr not match with expected")
	}
}

func TestEncodeInvalid(t *testing.T) {

}

func TestDecodeValid(t *testing.T) {

}

func TestDecodeInvalid(t *testing.T) {

}
