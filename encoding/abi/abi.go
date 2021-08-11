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

	// only can be applied to `uint` size <N> or `ufixed` size <N>
	unsignedTypeSize uint16
	// only can be applied to `ufixed` precision <M>
	unsignedTypePrecision uint16

	// length for static array / tuple
	staticLength uint32
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
		return t.childTypes[0].String() + "[" + strconv.Itoa(int(t.staticLength)) + "]"
	case Address:
		return "address"
	case ArrayDynamic:
		return t.childTypes[0].String() + "[]"
	case String:
		return "string"
	case Tuple:
		typeStrings := make([]string, len(t.childTypes))
		for i := 0; i < len(t.childTypes); i++ {
			typeStrings[i] = t.childTypes[i].String()
		}
		return "(" + strings.Join(typeStrings, ",") + ")"
	default:
		return "Bruh you should not be here"
	}
}

// TypeFromString de-serialization
func TypeFromString(str string) (Type, error) {
	switch {
	case strings.HasPrefix(str, "uint"):
		typeSize, err := strconv.ParseUint(str[4:], 10, 16)
		if err != nil {
			return Type{}, fmt.Errorf("ill formed uint type: %s", str)
		}
		uintTypeRes, err := MakeUintType(uint16(typeSize))
		if err != nil {
			return Type{}, err
		}
		return uintTypeRes, nil
	case str == "byte":
		return MakeByteType(), nil
	case strings.HasPrefix(str, "ufixed"):
		stringMatches := regexp.MustCompile(`^ufixed([\d]+)x([\d]+)`).FindStringSubmatch(str)
		// match string itself, then type-size, and type-precision
		if len(stringMatches) != 3 {
			return Type{}, fmt.Errorf("ill formed ufixed type: %s", str)
		}
		// guaranteed that there are 2 uint strings in ufixed string
		ufixedSize, err := strconv.ParseUint(stringMatches[1], 10, 16)
		if err != nil {
			return Type{}, err
		}
		ufixedPrecision, err := strconv.ParseUint(stringMatches[2], 10, 16)
		if err != nil {
			return Type{}, err
		}
		ufixedTypeRes, err := MakeUFixedType(uint16(ufixedSize), uint16(ufixedPrecision))
		if err != nil {
			return Type{}, err
		}
		return ufixedTypeRes, nil
	case str == "bool":
		return MakeBoolType(), nil
	case strings.HasPrefix(str, "]") && unicode.IsDigit(rune(str[len(str)-2])):
		stringMatches := regexp.MustCompile(`^[a-z\d\[\](),]+\[([\d]+)]$`).FindStringSubmatch(str)
		// match the string itself, then array length
		if len(stringMatches) != 2 {
			return Type{}, fmt.Errorf("static array ill formated: %s", str)
		}
		// guaranteed that the length of array is existing
		arrayLengthStr := stringMatches[1]
		arrayLength, err := strconv.ParseUint(arrayLengthStr, 10, 32)
		if err != nil {
			return Type{}, err
		}
		// parse the array element type
		arrayType, err := TypeFromString(str[:len(str)-(2+len(arrayLengthStr))])
		if err != nil {
			return Type{}, err
		}
		return MakeStaticArrayType(arrayType, uint32(arrayLength)), nil
	case str == "address":
		return MakeAddressType(), nil
	case strings.HasSuffix(str, "[]"):
		arrayArgType, err := TypeFromString(str[:len(str)-2])
		if err != nil {
			return arrayArgType, err
		}
		return MakeDynamicArrayType(arrayArgType), nil
	case str == "string":
		return MakeStringType(), nil
	case len(str) > 2 && str[0] == '(' && str[len(str)-1] == ')':
		tupleContent, err := parseTupleContent(strings.TrimSpace(str[1 : len(str)-1]))
		if err != nil {
			return Type{}, err
		} else if len(tupleContent) == 0 {
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
		return Type{}, fmt.Errorf("cannot convert a string %s to an ABI type", str)
	}
}

func parseTupleContent(str string) ([]string, error) {
	type segmentIndex struct{ left, right int }

	parenSegmentRecord, stack := make([]segmentIndex, 0), make([]int, 0)
	// get the most exterior parentheses segment (not overlapped by other parentheses)
	for index, chr := range str {
		if chr == '(' {
			stack = append(stack, index)
		} else if chr == ')' {
			if len(stack) == 0 {
				return []string{}, fmt.Errorf("unpaired parentheses: %s", str)
			}
			leftParenIndex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				parenSegmentRecord = append(parenSegmentRecord, segmentIndex{
					left:  leftParenIndex,
					right: index,
				})
			}
		}
	}
	if len(stack) != 0 {
		return []string{}, fmt.Errorf("unpaired parentheses: %s", str)
	}

	segmentRecord := make([]segmentIndex, 0)
	// iterate through parentheses segment and separate string into segments
	for _, seg := range parenSegmentRecord {
		if len(segmentRecord) == 0 {
			if seg.left != 0 {
				segmentRecord = append(segmentRecord, segmentIndex{
					left:  0,
					right: seg.left - 1,
				})
			}
		} else {
			prevRight := segmentRecord[len(segmentRecord)-1].right
			if prevRight+1 < seg.left {
				segmentRecord = append(segmentRecord, segmentIndex{
					left:  prevRight + 1,
					right: seg.left - 1,
				})
			}
		}
		segmentRecord = append(segmentRecord, seg)
	}
	// last segment, or only 1 segment case
	if len(segmentRecord) > 0 {
		prevRight := segmentRecord[len(segmentRecord)-1].right
		if prevRight != len(str)-1 {
			segmentRecord = append(segmentRecord, segmentIndex{
				left:  prevRight + 1,
				right: len(str) - 1,
			})
		}
	} else {
		segmentRecord = append(segmentRecord, segmentIndex{left: 0, right: len(str) - 1})
	}

	tupleContent := make([]string, 0)
	for _, segment := range segmentRecord {
		segmentStr := str[segment.left : segment.right+1]
		segmentStr = strings.Trim(segmentStr, ",")
		if strings.HasPrefix(segmentStr, "(") {
			tupleContent = append(tupleContent, segmentStr)
		} else {
			segmentStrs := strings.Split(segmentStr, ",")
			tupleContent = append(tupleContent, segmentStrs...)
		}
	}
	return tupleContent, nil
}

func MakeUintType(typeSize uint16) (Type, error) {
	if typeSize%8 != 0 || typeSize < 8 || typeSize > 512 {
		return Type{}, fmt.Errorf("type uint size mod 8 = 0, range [8, 512], error typesize: %d", typeSize)
	}
	return Type{
		typeFromEnum:     Uint,
		unsignedTypeSize: typeSize,
	}, nil
}

func MakeByteType() Type {
	return Type{
		typeFromEnum: Byte,
	}
}

func MakeUFixedType(typeSize uint16, typePrecision uint16) (Type, error) {
	if typeSize%8 != 0 || typeSize < 8 || typeSize > 512 {
		return Type{}, fmt.Errorf("type uint size mod 8 = 0, range [8, 512], error typesize: %d", typeSize)
	}
	if typePrecision > 160 || typePrecision < 1 {
		return Type{}, fmt.Errorf("type uint precision range [1, 160]")
	}
	return Type{
		typeFromEnum:          Ufixed,
		unsignedTypeSize:      typeSize,
		unsignedTypePrecision: typePrecision,
	}, nil
}

func MakeBoolType() Type {
	return Type{
		typeFromEnum: Bool,
	}
}

func MakeStaticArrayType(argumentType Type, arrayLength uint32) Type {
	return Type{
		typeFromEnum: ArrayStatic,
		childTypes:   []Type{argumentType},
		staticLength: arrayLength,
	}
}

func MakeAddressType() Type {
	return Type{
		typeFromEnum: Address,
	}
}

func MakeDynamicArrayType(argumentType Type) Type {
	return Type{
		typeFromEnum: ArrayDynamic,
		childTypes:   []Type{argumentType},
	}
}

func MakeStringType() Type {
	return Type{
		typeFromEnum: String,
	}
}

func MakeTupleType(argumentTypes []Type) Type {
	return Type{
		typeFromEnum: Tuple,
		childTypes:   argumentTypes,
		staticLength: uint32(len(argumentTypes)),
	}
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
	return []byte{}
}

// Decode de-serialization
func Decode(valueByte []byte, valueType Type) (Value, error) {
	// TODO
	return Value{}, nil
}

func (v Value) getUint() (uint64, error) {
	// TODO: how to handle different integer sizes and precision?

	// if v.valueType is not a uint64, return error
	return uint64(0), nil
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
	// change the value type later
}

func MakeValueUfixed(value *big.Int, size uint16, precision uint16) Value {
	// TODO: also consider how to handle differnet int sizes and precision
	return Value{}
}
