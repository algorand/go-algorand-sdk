package abi

import (
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
