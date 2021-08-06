package abi

import "math/big"

// All ABI types
/*
uint<N>: An N-bit unsigned integer, where 8 <= N <= 512 and N % 8 = 0.
byte // alias for uint8
ufixed<N>x<M>
bool
<type>[<N>]
address // alias for byte[32]
<type>[]
string
(T1,T2,...,TN):
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
}

// String serialization
func (t Type) String() string {
	// TODO
}

// TypeFromstring de-serialization
func TypeFromString(str string) (Type, error) {
	// TODO
}

// TODO makeType

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
