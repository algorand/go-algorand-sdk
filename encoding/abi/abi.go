package abi

import (
	"github.com/algorand/go-algorand/data/abi"
)

type Type = abi.Type

func TypeOf(str string) (Type, error) {
	return abi.TypeOf(str)
}
