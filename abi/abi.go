package abi

import (
	"fmt"
	"math"

	"github.com/algorand/go-algorand/data/abi"
)

type Type = abi.Type
type Value struct {
	AbiType  Type
	RawValue interface{}
}

func TypeOf(str string) (Type, error) {
	return abi.TypeOf(str)
}

// MakeTupleType makes tuple ABI type by taking an array of tuple element types as argument.
func MakeTupleType(argumentTypes []Type) (Type, error) {
	if len(argumentTypes) == 0 {
		return Type{}, fmt.Errorf("tuple must contain at least one type")
	}

	if len(argumentTypes) >= math.MaxUint16 {
		return Type{}, fmt.Errorf("tuple type child type number larger than maximum uint16 error")
	}

	strTuple := "(" + argumentTypes[0].String()
	for i := 1; i < len(argumentTypes); i++ {
		strTuple += "," + argumentTypes[i].String()
	}
	strTuple += ")"

	return TypeOf(strTuple)
}
