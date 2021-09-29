package abi

import (
	"math/big"

	"github.com/algorand/go-algorand/data/abi"
)

type ABIValue struct {
	abiVal abi.Value
}

type ABIType struct {
	abiType abi.Type
}

func (val ABIValue) Encode() ([]byte, error) {
	return val.abiVal.Encode()
}

func Decode(encoded []byte, abiType ABIType) (ABIValue, error) {
	abiVal, err := abi.Decode(encoded, abiType.abiType)
	if err != nil {
		return ABIValue{}, err
	}
	return ABIValue{abiVal}, nil
}

func (t ABIType) String() string {
	return t.abiType.String()
}

func TypeFromString(str string) (ABIType, error) {
	abiTypeFromStr, err := abi.TypeFromString(str)
	if err != nil {
		return ABIType{}, err
	}
	return ABIType{abiTypeFromStr}, nil
}

func MakeUint8(val uint8) ABIValue {
	return ABIValue{abi.MakeUint8(val)}
}

func (val ABIValue) GetUint8() (uint8, error) {
	return val.abiVal.GetUint8()
}

func MakeUint16(val uint16) ABIValue {
	return ABIValue{abi.MakeUint16(val)}
}

func (val ABIValue) GetUint16() (uint16, error) {
	return val.abiVal.GetUint16()
}

func MakeUint32(val uint32) ABIValue {
	return ABIValue{abi.MakeUint32(val)}
}

func (val ABIValue) GetUint32() (uint32, error) {
	return val.abiVal.GetUint32()
}

func MakeUint64(val uint64) ABIValue {
	return ABIValue{abi.MakeUint64(val)}
}

func (val ABIValue) GetUint64() (uint64, error) {
	return val.abiVal.GetUint64()
}

func MakeUint(val *big.Int, size uint16) (ABIValue, error) {
	uintV, err := abi.MakeUint(val, size)
	if err != nil {
		return ABIValue{}, err
	}
	return ABIValue{uintV}, nil
}

func (val ABIValue) GetUint() (*big.Int, error) {
	return val.abiVal.GetUint()
}

func MakeUfixed(val *big.Int, size uint16, precision uint16) (ABIValue, error) {
	ufixedV, err := abi.MakeUfixed(val, size, precision)
	if err != nil {
		return ABIValue{}, err
	}
	return ABIValue{ufixedV}, nil
}

func (val ABIValue) GetUfixed() (*big.Int, error) {
	return val.abiVal.GetUfixed()
}

func MakeString(val string) ABIValue {
	return ABIValue{abi.MakeString(val)}
}

func (val ABIValue) GetString() (string, error) {
	return val.abiVal.GetString()
}

func MakeByte(val byte) ABIValue {
	return ABIValue{abi.MakeByte(val)}
}

func (val ABIValue) GetByte() (byte, error) {
	return val.abiVal.GetByte()
}

func MakeAddress(val [32]byte) ABIValue {
	return ABIValue{abi.MakeAddress(val)}
}

func (val ABIValue) GetAddress() ([32]byte, error) {
	return val.abiVal.GetAddress()
}

func MakeBool(val bool) ABIValue {
	return ABIValue{abi.MakeBool(val)}
}

func (val ABIValue) GetBool() (bool, error) {
	return val.abiVal.GetBool()
}

func MakeDynamicArray(values []ABIValue, elemType ABIType) (ABIValue, error) {
	tempValues := make([]abi.Value, len(values))
	for i := 0; i < len(values); i++ {
		tempValues[i] = values[i].abiVal
	}
	abiDynamicArray, err := abi.MakeDynamicArray(tempValues, elemType.abiType)
	if err != nil {
		return ABIValue{}, err
	}
	return ABIValue{abiDynamicArray}, nil
}

func MakeStaticArray(values []ABIValue) (ABIValue, error) {
	tempValues := make([]abi.Value, len(values))
	for i := 0; i < len(values); i++ {
		tempValues[i] = values[i].abiVal
	}
	abiStaticArray, err := abi.MakeStaticArray(tempValues)
	if err != nil {
		return ABIValue{}, err
	}
	return ABIValue{abiStaticArray}, nil
}

func MakeTuple(values []ABIValue) (ABIValue, error) {
	tempValues := make([]abi.Value, len(values))
	for i := 0; i < len(values); i++ {
		tempValues[i] = values[i].abiVal
	}
	abiTuple, err := abi.MakeTuple(tempValues)
	if err != nil {
		return ABIValue{}, err
	}
	return ABIValue{abiTuple}, nil
}

func (val ABIValue) GetValueByIndex(index uint16) (ABIValue, error) {
	abiIndexVal, err := val.abiVal.GetValueByIndex(index)
	if err != nil {
		return ABIValue{}, err
	}
	return ABIValue{abiIndexVal}, nil
}

func MakeUintType(size uint16) (ABIType, error) {
	abiUintT, err := abi.MakeUintType(size)
	if err != nil {
		return ABIType{}, err
	}
	return ABIType{abiUintT}, nil
}

func MakeUfixedType(size uint16, precision uint16) (ABIType, error) {
	abiUfixedT, err := abi.MakeUfixedType(size, precision)
	if err != nil {
		return ABIType{}, err
	}
	return ABIType{abiUfixedT}, nil
}

func MakeByteType() ABIType {
	return ABIType{abi.MakeByteType()}
}

func MakeAddressType() ABIType {
	return ABIType{abi.MakeAddressType()}
}

func MakeStringType() ABIType {
	return ABIType{abi.MakeStringType()}
}

func MakeBoolType() ABIType {
	return ABIType{abi.MakeBoolType()}
}

func MakeDynamicArrayType(elemType ABIType) ABIType {
	return ABIType{abi.MakeDynamicArrayType(elemType.abiType)}
}

func MakeStaticArrayType(elemType ABIType, length uint16) ABIType {
	return ABIType{abi.MakeStaticArrayType(elemType.abiType, length)}
}

func MakeTupleType(elemTypes []ABIType) (ABIType, error) {
	tempTypes := make([]abi.Type, len(elemTypes))
	for i := 0; i < len(elemTypes); i++ {
		tempTypes[i] = elemTypes[i].abiType
	}
	tupleT, err := abi.MakeTupleType(tempTypes)
	if err != nil {
		return ABIType{}, err
	}
	return ABIType{tupleT}, nil
}
