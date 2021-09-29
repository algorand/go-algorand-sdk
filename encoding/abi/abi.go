package abi

import (
	"math/big"

	"github.com/algorand/go-algorand/data/abi"
)

type Value struct {
	abi.Value
}

type Type struct {
	abi.Type
}

func (val Value) Encode() ([]byte, error) {
	return val.Value.Encode()
}

func Decode(encoded []byte, abiType Type) (Value, error) {
	abiVal, err := abi.Decode(encoded, abiType.Type)
	if err != nil {
		return Value{}, err
	}
	return Value{abiVal}, nil
}

func (t Type) String() string {
	return t.Type.String()
}

func TypeFromString(str string) (Type, error) {
	abiTypeFromStr, err := abi.TypeFromString(str)
	if err != nil {
		return Type{}, err
	}
	return Type{abiTypeFromStr}, nil
}

func MakeUint8(val uint8) Value {
	return Value{abi.MakeUint8(val)}
}

func (val Value) GetUint8() (uint8, error) {
	return val.Value.GetUint8()
}

func MakeUint16(val uint16) Value {
	return Value{abi.MakeUint16(val)}
}

func (val Value) GetUint16() (uint16, error) {
	return val.Value.GetUint16()
}

func MakeUint32(val uint32) Value {
	return Value{abi.MakeUint32(val)}
}

func (val Value) GetUint32() (uint32, error) {
	return val.Value.GetUint32()
}

func MakeUint64(val uint64) Value {
	return Value{abi.MakeUint64(val)}
}

func (val Value) GetUint64() (uint64, error) {
	return val.Value.GetUint64()
}

func MakeUint(val *big.Int, size uint16) (Value, error) {
	uintV, err := abi.MakeUint(val, size)
	if err != nil {
		return Value{}, err
	}
	return Value{uintV}, nil
}

func (val Value) GetUint() (*big.Int, error) {
	return val.Value.GetUint()
}

func MakeUfixed(val *big.Int, size uint16, precision uint16) (Value, error) {
	ufixedV, err := abi.MakeUfixed(val, size, precision)
	if err != nil {
		return Value{}, err
	}
	return Value{ufixedV}, nil
}

func (val Value) GetUfixed() (*big.Int, error) {
	return val.Value.GetUfixed()
}

func MakeString(val string) Value {
	return Value{abi.MakeString(val)}
}

func (val Value) GetString() (string, error) {
	return val.Value.GetString()
}

func MakeByte(val byte) Value {
	return Value{abi.MakeByte(val)}
}

func (val Value) GetByte() (byte, error) {
	return val.Value.GetByte()
}

func MakeAddress(val [32]byte) Value {
	return Value{abi.MakeAddress(val)}
}

func (val Value) GetAddress() ([32]byte, error) {
	return val.Value.GetAddress()
}

func MakeBool(val bool) Value {
	return Value{abi.MakeBool(val)}
}

func (val Value) GetBool() (bool, error) {
	return val.Value.GetBool()
}

func MakeDynamicArray(values []Value, elemType Type) (Value, error) {
	tempValues := make([]abi.Value, len(values))
	for i := 0; i < len(values); i++ {
		tempValues[i] = values[i].Value
	}
	abiDynamicArray, err := abi.MakeDynamicArray(tempValues, elemType.Type)
	if err != nil {
		return Value{}, err
	}
	return Value{abiDynamicArray}, nil
}

func MakeStaticArray(values []Value) (Value, error) {
	tempValues := make([]abi.Value, len(values))
	for i := 0; i < len(values); i++ {
		tempValues[i] = values[i].Value
	}
	abiStaticArray, err := abi.MakeStaticArray(tempValues)
	if err != nil {
		return Value{}, err
	}
	return Value{abiStaticArray}, nil
}

func MakeTuple(values []Value) (Value, error) {
	tempValues := make([]abi.Value, len(values))
	for i := 0; i < len(values); i++ {
		tempValues[i] = values[i].Value
	}
	abiTuple, err := abi.MakeTuple(tempValues)
	if err != nil {
		return Value{}, err
	}
	return Value{abiTuple}, nil
}

func (val Value) GetValueByIndex(index uint16) (Value, error) {
	abiIndexVal, err := val.Value.GetValueByIndex(index)
	if err != nil {
		return Value{}, err
	}
	return Value{abiIndexVal}, nil
}

func MakeUintType(size uint16) (Type, error) {
	abiUintT, err := abi.MakeUintType(size)
	if err != nil {
		return Type{}, err
	}
	return Type{abiUintT}, nil
}

func MakeUfixedType(size uint16, precision uint16) (Type, error) {
	abiUfixedT, err := abi.MakeUfixedType(size, precision)
	if err != nil {
		return Type{}, err
	}
	return Type{abiUfixedT}, nil
}

func MakeByteType() Type {
	return Type{abi.MakeByteType()}
}

func MakeAddressType() Type {
	return Type{abi.MakeAddressType()}
}

func MakeStringType() Type {
	return Type{abi.MakeStringType()}
}

func MakeBoolType() Type {
	return Type{abi.MakeBoolType()}
}

func MakeDynamicArrayType(elemType Type) Type {
	return Type{abi.MakeDynamicArrayType(elemType.Type)}
}

func MakeStaticArrayType(elemType Type, length uint16) Type {
	return Type{abi.MakeStaticArrayType(elemType.Type, length)}
}

func MakeTupleType(elemTypes []Type) (Type, error) {
	tempTypes := make([]abi.Type, len(elemTypes))
	for i := 0; i < len(elemTypes); i++ {
		tempTypes[i] = elemTypes[i].Type
	}
	tupleT, err := abi.MakeTupleType(tempTypes)
	if err != nil {
		return Type{}, err
	}
	return Type{tupleT}, nil
}
