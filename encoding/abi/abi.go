package abi

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

/*
   ABI-Types: uint<N>: An N-bit unsigned integer (8 <= N <= 512 and N % 8 = 0).
            | byte (alias for uint8)
            | ufixed <N> x <M> (8 <= N <= 512, N % 8 = 0, and 0 < M <= 160)
            | bool
            | address (alias for byte[32])
            | <type> [<N>]
            | <type> []
            | string
            | (T1, ..., Tn)
*/

type BaseType uint32

const (
	Uint BaseType = iota
	Byte
	Ufixed
	Bool
	ArrayStatic
	Address
	ArrayDynamic
	String
	Tuple
)

// Need a struct which represents an ABI type. In the case of tuples and arrays,
// a type may have "children" types as well
// - Need a method to convert this to a string, e.g. (uint64,byte[]) for a tuple
// - Need a function to parse such a string back to this struct

type Type struct {
	typeFromEnum BaseType
	childTypes   []Type

	// only appliable to `uint` size <N> or `ufixed` size <N>
	unsignedTypeSize uint16
	// only appliable to `ufixed` precision <M>
	unsignedTypePrecision uint16

	// length for static array
	staticArrayLength uint64
}

// String serialization
func (t Type) String() string {
	switch t.typeFromEnum {
	case Uint:
		return "uint" + strconv.Itoa(int(t.unsignedTypeSize))
	case Byte:
		return "byte"
	case Ufixed:
		return "ufixed" + strconv.Itoa(int(t.unsignedTypeSize)) + "x" + strconv.Itoa(int(t.unsignedTypePrecision))
	case Bool:
		return "bool"
	case ArrayStatic:
		return "[" + strconv.FormatUint(t.staticArrayLength, 10) + "]" + t.childTypes[0].String()
	case Address:
		return "address"
	case ArrayDynamic:
		return "[]" + t.childTypes[0].String()
	case String:
		return "string"
	case Tuple:
		typeStrings := make([]string, len(t.childTypes))
		for i := 0; i < len(t.childTypes); i++ {
			typeStrings[i] = t.childTypes[i].String()
		}
		return "(" + strings.Join(typeStrings, ", ") + ")"
	default:
		return "Bruh you should not be here"
	}
}

// TypeFromString de-serialization
func TypeFromString(str string) (Type, error) {
	trimmedStr := strings.TrimSpace(str)
	switch {
	case len(trimmedStr) > 4 && trimmedStr[:4] == "uint":
		typeSize, err := strconv.ParseUint(trimmedStr[4:], 10, 16)
		if err != nil {
			// uint + not-uint value appended
			return Type{}, fmt.Errorf("ill formed uint type: %s", trimmedStr)
		}
		if typeSize % 8 != 0 || typeSize < 8 || typeSize > 512 {
			// uint + uint invalid value case
			return MakeAddressType(),
				fmt.Errorf("type uint size mod 8 = 0, range [8, 512], error type: %s", trimmedStr)
		}
		return MakeUintType(uint16(typeSize)), nil
	case trimmedStr == "byte":
		return MakeByteType(), nil
	case len(trimmedStr) > 6 && trimmedStr[:6] == "ufixed":
		match, err := regexp.MatchString(trimmedStr, `^ufixed[\d]+x[\d]+$`)
		if err != nil {
			return Type{}, err
		}
		if !match {
			return Type{}, fmt.Errorf("ufixed type ill formated: %s", trimmedStr)
		}
		re := regexp.MustCompile(`[\d]+`)
		// guaranteed that there are 2 uint strings in ufixed string
		ufixedNums := re.FindAllString(trimmedStr[6:], 2)
		ufixedSize, err := strconv.ParseUint(ufixedNums[0], 10, 16)
		if err != nil {
			return Type{}, err
		}
		ufixedPrecision, err := strconv.ParseUint(ufixedNums[1], 10, 16)
		if err != nil {
			return Type{}, err
		}
		return MakeUFixedType(uint16(ufixedSize), uint16(ufixedPrecision)), nil
	case trimmedStr == "bool":
		return MakeBoolType(), nil
	case len(trimmedStr) > 2 && trimmedStr[0] == '[' && unicode.IsDigit(rune(trimmedStr[1])):
		match, err := regexp.MatchString(trimmedStr, `^\[[\d]+\].+$`)
		if err != nil {
			return Type{}, err
		}
		if !match {
			return Type{}, fmt.Errorf("static array ill formated: %s", trimmedStr)
		}
		re := regexp.MustCompile(`[\d]+`)
		// guaranteed that the length of array is existing
		arrayLengthStrArray := re.FindAllString(trimmedStr, 1)
		arrayLength, err := strconv.ParseUint(arrayLengthStrArray[0], 10, 64)
		if err != nil {
			return Type{}, err
		}
		// parse the array element type
		arrayType, err := TypeFromString(trimmedStr[2 + len(arrayLengthStrArray[0]):])
		if err != nil {
			return Type{}, err
		}
		return MakeStaticArrayType(arrayType, arrayLength), nil
	case trimmedStr == "address":
		return MakeAddressType(), nil
	case len(trimmedStr) > 2 && trimmedStr[:2] == "[]":
		arrayArgType, err := TypeFromString(trimmedStr[2:])
		if err != nil {
			return arrayArgType, err
		}
		return MakeDynamicArrayType(arrayArgType), nil
	case trimmedStr == "string":
		return MakeStringType(), nil
	case len(trimmedStr) > 2 && trimmedStr[0] == '(' && trimmedStr[len(trimmedStr) - 1] == ')':
		tupleContent := strings.Split(strings.TrimSpace(trimmedStr[1: len(trimmedStr) - 1]), ",")
		if len(tupleContent) == 0 {
			return Type{}, fmt.Errorf("tuple type has no argument types")
		}
		tupleTypes := make([]Type, len(tupleContent))
		for i := 0; i < len(tupleContent); i++ {
			ti, err := TypeFromString(tupleContent[i])
			if err != nil {
				return Type{}, err
			}
			tupleTypes[i] = ti
		}
		return MakeTupleType(tupleTypes), nil
	default:
		return Type{}, fmt.Errorf("cannot convert a string %s to an ABI type", trimmedStr)
	}
}

func MakeUintType(typeSize uint16) Type {
	return Type{Uint, nil, typeSize, 0, 0}
}

func MakeByteType() Type {
	return Type{Byte, nil, 0, 0, 0}
}

func MakeUFixedType(typeSize uint16, typePrecision uint16) Type {
	return Type{Ufixed, nil, typeSize, typePrecision, 0}
}

func MakeBoolType() Type {
	return Type{Bool, nil, 0, 0, 0}
}

func MakeStaticArrayType(argumentType Type, arrayLength uint64) Type {
	return Type{ArrayStatic, []Type{argumentType}, 0, 0, arrayLength}
}

func MakeAddressType() Type {
	return Type{Address, nil, 0, 0, 0}
}

func MakeDynamicArrayType(argumentType Type) Type {
	return Type{ArrayDynamic, []Type{argumentType}, 0, 0, 0}
}

func MakeStringType() Type {
	return Type{String, nil, 0, 0, 0}
}

func MakeTupleType(argumentTypes []Type) Type {
	return Type{Tuple, argumentTypes, 0, 0, uint64(len(argumentTypes))}
}

// Need a struct which represents an ABI value. This struct would probably
// contain the ABI type struct and an interface{} for its value(s)
// - Need a way to encode this struct to bytes
// - Need a way to decode bytes into this struct
// - Need a way for the user to create and populate this struct
//    - This encompasses two cases: creating Values from golang values, and
//		modifying existing Value structs. For now creating is more important.
// - Need a way for the user to get specific values from this struct

type Value struct {
	valueType Type
	value     interface{}
}

// Encode serialization
func (v Value) Encode() []byte {
	// TODO
}

// Decode de-serialization
func Decode(valueByte []byte, valueType Type) (Value, error) {
	// TODO
}

func (v Value) getUint() (uint64, error) {
	// TODO: how to handle different integer sizes and precision?

	// if v.valueType is not a uint64, return error
}

// TODO create get... functions
// TODO if its get array/tuple function, pass index and take the element

func MakeValueUint8(value uint8) Value {
	bigInt := big.NewInt(int64(value))
	return MakeValueUint(bigInt, 8)
}

func MakeValueUint16(value uint16) Value {
	bigInt := big.NewInt(int64(value))
	return MakeValueUint(bigInt, 16)
}

func MakeValueUint32(value uint32) Value {
	bigInt := big.NewInt(int64(value))
	return MakeValueUint(bigInt, 32)
}

func MakeValueUint64(value uint64) Value {
	bigInt := big.NewInt(int64(0))
	bigInt.SetUint64(value)
	return MakeValueUint(bigInt, 64)
}

func MakeValueUint(value *big.Int, size uint16) Value {
	return MakeValueUfixed(value, size, 0)
}

func MakeValueUfixed(value *big.Int, size uint16, precision uint16) Value {
	// TODO: also consider how to handle differnet int sizes and precision
}
