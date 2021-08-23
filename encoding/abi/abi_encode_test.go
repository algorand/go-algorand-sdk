package abi

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"testing"

	"github.com/chrismcguire/gobberish"
	"github.com/stretchr/testify/require"
)

// test with go test -cover -v

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
			randomIntByte := randomInt.Bytes()
			buffer := make([]byte, intSize/8-len(randomIntByte))
			buffer = append(buffer, randomIntByte...)
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

				randomBytes := randomInt.Bytes()
				buffer := make([]byte, size/8-len(randomBytes))
				buffer = append(buffer, randomBytes...)
				require.Equal(t, buffer, encodedUfixed, "encode ufixed not match with expected")
			}
		}
	}

	upperLimit := big.NewInt(0).Lsh(big.NewInt(1), 256)
	for i := 0; i < 1000; i++ {
		randomAddrInt, err := rand.Int(rand.Reader, upperLimit)
		require.NoError(t, err, "cryptographic random int init fail")

		addressBytes := randomAddrInt.Bytes()
		address := make([]byte, 32-len(addressBytes))
		address = append(address, addressBytes...)

		var addrBytes [32]byte
		copy(addrBytes[:], address[:32])

		addressValue := MakeAddress(addrBytes)
		addrEncode, err := addressValue.Encode()
		require.NoError(t, err, "address encode fail")
		require.Equal(t, address, addrEncode, "encode addr not match with expected")
	}

	for i := 0; i < 2; i++ {
		boolValue := MakeBool(i == 1)
		boolEncode, err := boolValue.Encode()
		require.NoError(t, err, "bool encode fail")
		expected := []byte{0x00}
		if i == 1 {
			expected = []byte{0x80}
		}
		require.Equal(t, expected, boolEncode, "encode bool not match with expected")
	}

	for i := 0; i < (1 << 8); i++ {
		byteValue := MakeByte(byte(i))
		byteEncode, err := byteValue.Encode()
		require.NoError(t, err, "byte encode fail")
		expected := []byte{byte(i)}
		require.Equal(t, expected, byteEncode, "encode byte not match with expected")
	}

	for length := 1; length <= 10; length++ {
		for i := 0; i < 10; i++ {
			utf8Str := gobberish.GenerateString(length)
			strValue := MakeString(utf8Str)
			utf8ByteLen := len([]byte(utf8Str))
			head := make([]byte, 2)
			binary.BigEndian.PutUint16(head, uint16(utf8ByteLen))
			expected := append(head, []byte(utf8Str)...)
			strEncode, err := strValue.Encode()
			require.NoError(t, err, "string encode fail")
			require.Equal(t, expected, strEncode, "encode string not match with expected")
		}
	}

	var testcases = []struct {
		input          []Value
		childT         []Type
		testType       string
		valueConstruct func([]Value, interface{}) (Value, error)
		expected       []byte
	}{
		// array static

		// array dynamic

		// tuple
	}

	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("test %s", testcase.testType), func(t *testing.T) {

		})
	}
}

func TestEncodeInvalid(t *testing.T) {

}

func TestDecodeValid(t *testing.T) {

}

func TestDecodeInvalid(t *testing.T) {

}
