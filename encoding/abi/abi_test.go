package abi

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestMakeUintTypeValid(t *testing.T) {
	for i := 8; i <= 512; i += 8 {
		uintType, _ := MakeUintType(uint16(i))
		expected := "uint" + strconv.Itoa(i)
		actual := uintType.String()
		require.Equal(t, expected, actual, "MakeUintType: expected %s, actual %s", expected, actual)
	}
}

func TestMakeUintTypeInvalid(t *testing.T) {
	// TODO something need to be added
	// uintType, err := MakeUintType(uint16(7))
	// require.Error(t, err)
	// require.Equal(t, "expected message", err.Error())
}

func TestTypeFromStringUintTypeValid(t *testing.T) {
	for i := 8; i <= 512; i += 8 {
		expected, _ := MakeUintType(uint16(i))
		actual, err := TypeFromString(expected.String())
		require.Equal(t, nil, err, "TypeFromString: uint parsing error: %s", expected.String())
		require.Equal(t, expected, actual, "TypeFromString: expected %s, actual %s", expected.String(), actual.String())
	}
}

func TestTypeFromStringUintTypeInvalid(t *testing.T) {
	// TODO something need to be added

	// testCases := []struct {
	// 	input string
	// 	expected Type
	// }{
	// 	{
	// 		input: "uint64",
	// 		expected: MakeUintType(uint16(64)),
	// 	},
	// 	{
	// 		input: "(uint64, (uint128, ufixed512x100))",
	// 		expected: MakeUintType(uint16(64)),
	// 	}
	// }

	// for index, test := range ... {
	// 	t.Run(fmt.Sprintf("test %d", index), func(t *testing.T) {
	// 		t.Error("abc")
	// 	})
	// }

}

func TestMakeUfixedTypeValid(t *testing.T) {
	for i := 8; i <= 512; i += 8 {
		for j := 1; j <= 160; j++ {
			ufixedType, _ := MakeUFixedType(uint16(i), uint16(j))
			expected := "ufixed" + strconv.Itoa(i) + "x" + strconv.Itoa(j)
			actual := ufixedType.String()
			require.Equal(t, expected, actual, "TypeFromString ufixed error: expected %s, actual %s", expected, actual)
		}
	}
}

func TestTypeFromStringUfixedTypeValid(t *testing.T) {
	for i := 8; i <= 512; i += 8 {
		for j := 1; j <= 160; j++ {
			expected, _ := MakeUFixedType(uint16(i), uint16(j))
			actual, err := TypeFromString("ufixed" + strconv.Itoa(i) + "x" + strconv.Itoa(j))
			require.Equal(t, nil, err, "TypeFromString ufixed parsing error: %s", expected.String())
			require.Equal(t, expected, actual, "TypeFromString ufixed: expected %s, actual %s", expected.String(), actual.String())
		}
	}
}

func TestMakeSimpleTypeValid(t *testing.T) {
	var testcases = []struct{
		input Type
		testType string
		expected string
	}{
		{input: MakeBoolType(), testType: "bool", expected: "bool"},
		{input: MakeStringType(), testType: "string", expected: "string"},
		{input: MakeAddressType(), testType: "address", expected: "address"},
		{input: MakeByteType(), testType: "byte", expected: "byte"},
	}
	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("MakeType test %s", testcase.testType), func(t *testing.T) {
			actual := testcase.input.String()
			require.Equal(t, testcase.expected, actual, "MakeType: expected %s, actual %s", testcase.expected, actual)
		})
	}
}

func TestSimpleTypeFromStringValid(t *testing.T) {
	var testcases = []struct {
		input    string
		testType string
		expected Type
	}{
		{input: MakeBoolType().String(), testType: "bool", expected: MakeBoolType()},
		{input: MakeStringType().String(), testType: "string", expected: MakeStringType()},
		{input: MakeAddressType().String(), testType: "address", expected: MakeAddressType()},
		{input: MakeByteType().String(), testType: "byte", expected: MakeByteType()},
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
