package abi

import "github.com/algorand/go-algorand/data/abi"

func Encode(val abi.Value) ([]byte, error) {
	return val.Encode()
}

func Decode(encoded []byte, abiType abi.Type) (abi.Value, error) {
	return abi.Decode(encoded, abiType)
}
