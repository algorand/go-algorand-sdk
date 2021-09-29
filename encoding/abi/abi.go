package abi

import (
	"math/big"

	"github.com/algorand/go-algorand/data/abi"
)

type Value = abi.Value
type Type = abi.Type

func Decode(encoded []byte, abiType Type) (Value, error) {
	return abi.Decode(encoded, abiType)
}

func TypeFromString(str string) (Type, error) {
	return abi.TypeFromString(str)
}

func MakeUint8(val uint8) Value {
	return abi.MakeUint8(val)
}

func MakeUint16(val uint16) Value {
	return abi.MakeUint16(val)
}

func MakeUint32(val uint32) Value {
	return abi.MakeUint32(val)
}

func MakeUint64(val uint64) Value {
	return abi.MakeUint64(val)
}

func MakeUint(val *big.Int, size uint16) (Value, error) {
	return abi.MakeUint(val, size)
}

func MakeUfixed(val *big.Int, size uint16, precision uint16) (Value, error) {
	return abi.MakeUfixed(val, size, precision)
}

func MakeString(val string) Value {
	return abi.MakeString(val)
}

func MakeByte(val byte) Value {
	return abi.MakeByte(val)
}

func MakeAddress(val [32]byte) Value {
	return abi.MakeAddress(val)
}

func MakeBool(val bool) Value {
	return abi.MakeBool(val)
}

func MakeDynamicArray(values []Value, elemType Type) (Value, error) {
	return abi.MakeDynamicArray(values, elemType)
}

func MakeStaticArray(values []Value) (Value, error) {
	return abi.MakeStaticArray(values)
}

func MakeTuple(values []Value) (Value, error) {
	return abi.MakeTuple(values)
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
