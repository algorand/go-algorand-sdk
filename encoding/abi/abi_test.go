package abi

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO need a fuzz test for the parsing

func TestMakeTypeValid(t *testing.T) {
	// uint
	for i := 8; i <= 512; i += 8 {
		uintType, _ := MakeUintType(uint16(i))
		expected := "uint" + strconv.Itoa(i)
		actual := uintType.String()
		require.Equal(t, expected, actual, "MakeUintType: expected %s, actual %s", expected, actual)
	}
	// ufixed
	for i := 8; i <= 512; i += 8 {
		for j := 1; j <= 160; j++ {
			ufixedType, _ := MakeUFixedType(uint16(i), uint16(j))
			expected := "ufixed" + strconv.Itoa(i) + "x" + strconv.Itoa(j)
			actual := ufixedType.String()
			require.Equal(t, expected, actual,
				"TypeFromString ufixed error: expected %s, actual %s", expected, actual)
		}
	}
	// bool/strings/address/byte + dynamic/static array + tuple
	var testcases = []struct {
		input    Type
		testType string
		expected string
	}{
		{input: MakeBoolType(), testType: "bool", expected: "bool"},
		{input: MakeStringType(), testType: "string", expected: "string"},
		{input: MakeAddressType(), testType: "address", expected: "address"},
		{input: MakeByteType(), testType: "byte", expected: "byte"},
		// dynamic array
		{
			input: MakeDynamicArrayType(
				Type{
					typeFromEnum:     Uint,
					unsignedTypeSize: uint16(32),
				},
			),
			testType: "dynamic array",
			expected: "uint32[]",
		},
		{
			input: MakeDynamicArrayType(
				MakeDynamicArrayType(
					MakeByteType(),
				),
			),
			testType: "dynamic array",
			expected: "byte[][]",
		},
		{
			input: MakeStaticArrayType(
				Type{
					typeFromEnum:          Ufixed,
					unsignedTypeSize:      uint16(128),
					unsignedTypePrecision: uint16(10),
				},
				uint32(100),
			),
			testType: "dynamic array",
			expected: "ufixed128x10[100]",
		},
		{
			input: MakeStaticArrayType(
				MakeStaticArrayType(
					MakeBoolType(),
					uint32(128),
				),
				uint32(256),
			),
			testType: "dynamic array",
			expected: "bool[128][256]",
		},
		// tuple type
		{
			input: MakeTupleType(
				[]Type{
					{
						typeFromEnum:     Uint,
						unsignedTypeSize: uint16(32),
					},
					MakeTupleType(
						[]Type{
							MakeAddressType(),
							MakeByteType(),
							MakeStaticArrayType(MakeBoolType(), uint32(10)),
							MakeDynamicArrayType(
								Type{
									typeFromEnum:          Ufixed,
									unsignedTypeSize:      uint16(256),
									unsignedTypePrecision: uint16(10),
								},
							),
						},
					),
					MakeDynamicArrayType(MakeByteType()),
				},
			),
			testType: "tuple type",
			expected: "(uint32,(address,byte,bool[10],ufixed256x10[]),byte[])",
		},
	}
	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("MakeType test %s", testcase.testType), func(t *testing.T) {
			actual := testcase.input.String()
			require.Equal(t, testcase.expected, actual,
				"MakeType: expected %s, actual %s", testcase.expected, actual)
		})
	}
}

func TestMakeTypeInvalid(t *testing.T) {
	// uint
	for i := 0; i <= 1000; i++ {
		randInput := rand.Uint32()
		for randInput%8 == 0 && randInput <= 512 && randInput >= 8 {
			randInput = rand.Uint32()
		}
		// note: if a var mod 8 = 0 (or not) in uint32, then it should mod 8 = 0 (or not) in uint16.
		_, err := MakeUintType(uint16(randInput))
		require.Error(t, err, "MakeUintType: should throw error on size input %d", randInput)
	}
	// ufixed
	for i := 0; i <= 10000; i++ {
		randSize := rand.Uint64()
		for randSize%8 == 0 && randSize <= 512 && randSize >= 8 {
			randSize = rand.Uint64()
		}
		randPrecision := rand.Uint32()
		for randPrecision >= 1 && randPrecision <= 160 {
			randPrecision = rand.Uint32()
		}
		_, err := MakeUFixedType(uint16(randSize), uint16(randPrecision))
		require.Error(t, err, "MakeUintType: should throw error on size input %d", randSize)
	}
}

func TestTypeFromStringValid(t *testing.T) {
	// uint
	for i := 8; i <= 512; i += 8 {
		expected, _ := MakeUintType(uint16(i))
		actual, err := TypeFromString(expected.String())
		require.Equal(t, nil, err, "TypeFromString: uint parsing error: %s", expected.String())
		require.Equal(t, expected, actual,
			"TypeFromString: expected %s, actual %s", expected.String(), actual.String())
	}
	// ufixed
	for i := 8; i <= 512; i += 8 {
		for j := 1; j <= 160; j++ {
			expected, _ := MakeUFixedType(uint16(i), uint16(j))
			actual, err := TypeFromString("ufixed" + strconv.Itoa(i) + "x" + strconv.Itoa(j))
			require.Equal(t, nil, err, "TypeFromString ufixed parsing error: %s", expected.String())
			require.Equal(t, expected, actual,
				"TypeFromString ufixed: expected %s, actual %s", expected.String(), actual.String())
		}
	}
	var testcases = []struct {
		input    string
		testType string
		expected Type
	}{
		{input: MakeBoolType().String(), testType: "bool", expected: MakeBoolType()},
		{input: MakeStringType().String(), testType: "string", expected: MakeStringType()},
		{input: MakeAddressType().String(), testType: "address", expected: MakeAddressType()},
		{input: MakeByteType().String(), testType: "byte", expected: MakeByteType()},
		{
			input:    "uint256[]",
			testType: "dynamic array",
			expected: MakeDynamicArrayType(Type{typeFromEnum: Uint, unsignedTypeSize: 256}),
		},
		{
			input:    "ufixed256x64[]",
			testType: "dynamic array",
			expected: MakeDynamicArrayType(
				Type{
					typeFromEnum:          Ufixed,
					unsignedTypeSize:      256,
					unsignedTypePrecision: 64,
				},
			),
		},
		{
			input:    "byte[][][][]",
			testType: "dynamic array",
			expected: MakeDynamicArrayType(
				MakeDynamicArrayType(
					MakeDynamicArrayType(
						MakeDynamicArrayType(
							MakeByteType(),
						),
					),
				),
			),
		},
		// static array
		{
			input:    "address[100]",
			testType: "static array",
			expected: MakeStaticArrayType(
				MakeAddressType(),
				uint32(100),
			),
		},
		{
			input:    "uint64[][200]",
			testType: "static array",
			expected: MakeStaticArrayType(
				MakeDynamicArrayType(
					Type{typeFromEnum: Uint, unsignedTypeSize: uint16(64)},
				),
				uint32(200),
			),
		},
		// tuple type
		{
			input:    "(uint32,(address,byte,bool[10],ufixed256x10[]),byte[])",
			testType: "tuple type",
			expected: MakeTupleType(
				[]Type{
					{
						typeFromEnum:     Uint,
						unsignedTypeSize: uint16(32),
					},
					MakeTupleType(
						[]Type{
							MakeAddressType(),
							MakeByteType(),
							MakeStaticArrayType(MakeBoolType(), uint32(10)),
							MakeDynamicArrayType(
								Type{
									typeFromEnum:          Ufixed,
									unsignedTypeSize:      uint16(256),
									unsignedTypePrecision: uint16(10),
								},
							),
						},
					),
					MakeDynamicArrayType(MakeByteType()),
				},
			),
		},
		{
			input:    "(uint32,(address,byte,bool[10],(ufixed256x10[])))",
			testType: "tuple type",
			expected: MakeTupleType(
				[]Type{
					{
						typeFromEnum:     Uint,
						unsignedTypeSize: uint16(32),
					},
					MakeTupleType(
						[]Type{
							MakeAddressType(),
							MakeByteType(),
							MakeStaticArrayType(MakeBoolType(), uint32(10)),
							MakeTupleType(
								[]Type{
									MakeDynamicArrayType(
										Type{
											typeFromEnum:          Ufixed,
											unsignedTypeSize:      uint16(256),
											unsignedTypePrecision: uint16(10),
										},
									),
								},
							),
						},
					),
				},
			),
		},
		{
			input:    "((uint32),(address,(byte,bool[10],ufixed256x10[])))",
			testType: "tuple type",
			expected: MakeTupleType(
				[]Type{
					MakeTupleType(
						[]Type{
							{
								typeFromEnum:     Uint,
								unsignedTypeSize: uint16(32),
							},
						},
					),
					MakeTupleType(
						[]Type{
							MakeAddressType(),
							MakeTupleType(
								[]Type{
									MakeByteType(),
									MakeStaticArrayType(MakeBoolType(), uint32(10)),
									MakeDynamicArrayType(
										Type{
											typeFromEnum:          Ufixed,
											unsignedTypeSize:      uint16(256),
											unsignedTypePrecision: uint16(10),
										},
									),
								},
							),
						},
					),
				},
			),
		},
	}
	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("TypeFromString test %s", testcase.testType), func(t *testing.T) {
			actual, err := TypeFromString(testcase.input)
			require.Equal(t, nil, err, "TypeFromString %s parsing error", testcase.testType)
			require.Equal(t, testcase.expected, actual, "TestFromString %s: expected %s, actual %s",
				testcase.testType, testcase.expected.String(), actual.String())
		})
	}
}

func TestTypeFromStringInvalid(t *testing.T) {
	for i := 0; i <= 1000; i++ {
		randSize := rand.Uint64()
		for randSize%8 == 0 && randSize <= 512 && randSize >= 8 {
			randSize = rand.Uint64()
		}
		errorInput := "uint" + strconv.FormatUint(randSize, 10)
		_, err := TypeFromString(errorInput)
		require.Error(t, err, "MakeUintType: should throw error on size input %d", randSize)
	}
	for i := 0; i <= 10000; i++ {
		randSize := rand.Uint64()
		for randSize%8 == 0 && randSize <= 512 && randSize >= 8 {
			randSize = rand.Uint64()
		}
		randPrecision := rand.Uint64()
		for randPrecision >= 1 && randPrecision <= 160 {
			randPrecision = rand.Uint64()
		}
		errorInput := "ufixed" + strconv.FormatUint(randSize, 10) + "x" + strconv.FormatUint(randPrecision, 10)
		_, err := TypeFromString(errorInput)
		require.Error(t, err, "MakeUintType: should throw error on size input %d", randSize)
	}
	var testcases = []string{
		// uint
		"uint123x345",
		"uint 128",
		"uint8 ",
		"uint!8",
		"uint[32]",
		"uint-893",
		"uint#120\\",
		// ufixed
		"ufixed000000000016x0000010",
		"ufixed123x345",
		"ufixed 128 x 100",
		"ufixed64x10 ",
		"ufixed!8x2 ",
		"ufixed[32]x16",
		"ufixed-64x+100",
		"ufixed16x+12",
		// dynamic array
		"uint256 []",
		"byte[] ",
		"[][][]",
		"stuff[]",
		// static array
	}
	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("TypeFromString dynamic array test %s", testcase), func(t *testing.T) {
			_, err := TypeFromString(testcase)
			require.Error(t, err, "%s should throw error", testcase)
		})
	}
}
