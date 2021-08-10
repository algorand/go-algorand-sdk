package abi

import (
	"reflect"
	"strconv"
	"testing"
)

func testValidFramework(funcName string, expected interface{}, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s testing error: got %s, expect %s", funcName, actual, expected)
	}
}

func TestMakeUintTypeValid(t *testing.T) {
	for i := 8; i <= 512; i += 8 {
		uintType := MakeUintType(uint16(i))
		expected := "uint" + strconv.Itoa(i)
		actual := uintType.String()
		testValidFramework("MakeUintType", expected, actual, t)
	}
}

func TestMakeUintTypeInvalid(t *testing.T) {
	// TODO something need to be added
}

func TestTypeFromStringUintTypeValid(t *testing.T) {
	for i := 8; i <= 512; i += 8 {
		expected := MakeUintType(uint16(i))
		actual, err := TypeFromString(expected.String())
		if err != nil {
			t.Errorf("TypeFromString testing error: Parsing error for %s", expected.String())
		}
		testValidFramework("TypeFromString uint", expected, actual, t)
	}
}

func TestTypeFromStringUintTypeInvalid(t *testing.T) {
	// TODO something need to be added
}

func TestMakeByteTypeValid(t *testing.T) {
	byteType := MakeByteType()
	expected := byteType.String()
	actual := "byte"
	testValidFramework("MakeByteType", expected, actual, t)
}
