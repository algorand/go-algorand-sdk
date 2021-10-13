package abi

import (
	"github.com/algorand/go-algorand/data/abi"
)

type Type = abi.Type

func TypeFromString(str string) (Type, error) {
	return abi.TypeFromString(str)
}

func MakeUintType(size uint16) (Type, error) {
	return abi.MakeUintType(size)
}

func MakeUfixedType(size uint16, precision uint16) (Type, error) {
	return abi.MakeUfixedType(size, precision)
}

func MakeByteType() Type {
	return abi.MakeByteType()
}

func MakeAddressType() Type {
	return abi.MakeAddressType()
}

func MakeStringType() Type {
	return abi.MakeStringType()
}

func MakeBoolType() Type {
	return abi.MakeBoolType()
}

func MakeDynamicArrayType(elemType Type) Type {
	return abi.MakeDynamicArrayType(elemType)
}

func MakeStaticArrayType(elemType Type, length uint16) Type {
	return abi.MakeStaticArrayType(elemType, length)
}

func MakeTupleType(elemTypes []Type) (Type, error) {
	return abi.MakeTupleType(elemTypes)
}
