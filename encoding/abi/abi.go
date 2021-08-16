package abi

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
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
	/*
		by ABI spec, len over binary array returns number of bytes
		the type is uint16, which allows for only lenth in [0, 2^16 - 1]
		representation of static length can only be constrained in uint16 type
	*/
	staticLength uint16
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
	case strings.HasSuffix(str, "[]"):
		arrayArgType, err := TypeFromString(str[:len(str)-2])
		if err != nil {
			return arrayArgType, err
		}
		return MakeDynamicArrayType(arrayArgType), nil
	case strings.HasSuffix(str, "]") && len(str) >= 2 && unicode.IsDigit(rune(str[len(str)-2])):
		stringMatches := regexp.MustCompile(`^[a-z\d\[\](),]+\[([1-9][\d]*)]$`).FindStringSubmatch(str)
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
		return MakeStaticArrayType(arrayType, uint16(arrayLength)), nil
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
		stringMatches := regexp.MustCompile(`^ufixed([1-9][\d]*)x([1-9][\d]*)$`).FindStringSubmatch(str)
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
	case str == "address":
		return MakeAddressType(), nil
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
	// illustration: "*****,(*****),*****" => ["*****", "(*****)", "*****"]
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
		if len(segmentStr) == 0 {
			if segment.right+1-segment.left > 1 {
				return []string{}, fmt.Errorf("no consequtive commas")
			} else {
				continue
			}
		}
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

func MakeStaticArrayType(argumentType Type, arrayLength uint16) Type {
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
		staticLength: uint16(len(argumentTypes)),
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
	switch v.valueType.typeFromEnum {
	case Uint:
		bigIntValue, _ := GetUint(v)
		buffer := make([]byte, v.valueType.unsignedTypeSize/8)
		return bigIntValue.FillBytes(buffer)
	case Ufixed:
		ufixedValue, _ := GetUfixed(v)
		buffer := make([]byte, v.valueType.unsignedTypeSize/8)
		return ufixedValue.Num().FillBytes(buffer)
	case Bool:
		return []byte{}
	case Byte:
		bytesValue, _ := GetByte(v)
		return []byte{bytesValue}
	case ArrayStatic:
		staticArrayBytes := make([]byte, 0)
		for i := 0; i < int(v.valueType.staticLength); i++ {
			element, _ := GetStaticArrayByIndex(v, uint16(i))
			staticArrayBytes = append(staticArrayBytes, element.Encode()...)
		}
		return staticArrayBytes
	case Address:
		addressValue, _ := GetAddress(v)
		return addressValue[:]
	case ArrayDynamic:
		dynamicArrayBytes := make([]byte, 2)
		arrayLen := reflect.ValueOf(v.value).Len()
		binary.BigEndian.PutUint16(dynamicArrayBytes, uint16(arrayLen))
		for i := 0; i < arrayLen; i++ {
			element, _ := GetDynamicArrayByIndex(v, uint16(i))
			dynamicArrayBytes = append(dynamicArrayBytes, element.Encode()...)
		}
		return dynamicArrayBytes
	case String:
		// need to rework, like dynamic array
		stringValue, _ := GetString(v)
		length := uint16(len(stringValue))
		stringBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(stringBytes, length)
		// encode length + string to byte array
		stringBytes = append(stringBytes, []byte(stringValue)...)
		return stringBytes
	case Tuple:
		// TODO need to check tuple specs
		return []byte{}
	default:
		return []byte("bruh you should not be here in encoding")
	}
}

// Decode de-serialization
func Decode(valueByte []byte, valueType Type) (Value, error) {
	switch valueType.typeFromEnum {
	case Uint:
		uintValue := big.NewInt(0).SetBytes(valueByte)
		return MakeUint(uintValue, valueType.unsignedTypeSize)
	case Ufixed:
		ufixedNumerator := big.NewInt(0).SetBytes(valueByte)
		ufixedDenominator := big.NewInt(0).Exp(
			big.NewInt(10), big.NewInt(int64(valueType.unsignedTypePrecision)),
			nil,
		)
		ufixedValue := big.NewRat(1, 1).SetFrac(ufixedNumerator, ufixedDenominator)
		return MakeUfixed(ufixedValue, valueType.unsignedTypeSize, valueType.unsignedTypePrecision)
	case Bool:
		return Value{}, nil
	case Byte:
		if len(valueByte) != 1 {
			return Value{}, fmt.Errorf("byte should be length 1")
		}
		return MakeByte(valueByte[0]), nil
	case ArrayStatic:
		return Value{}, nil
	case Address:
		if len(valueByte) != 32 {
			return Value{}, fmt.Errorf("address should be length 32")
		}
		var byteAssign [32]byte
		copy(byteAssign[:], valueByte)
		return MakeAddress(byteAssign), nil
	case ArrayDynamic:
		return Value{}, nil
	case String:
		if len(valueByte) <= 2 {
			return Value{}, fmt.Errorf("string format corrupted")
		}
		stringLenBytes := valueByte[:2]
		stringLen := binary.BigEndian.Uint16(stringLenBytes)
		stringValue := string(valueByte[2:])
		if len(stringValue) != int(stringLen) {
			return Value{}, fmt.Errorf("string value do not match string length")
		}
		return MakeString(stringValue), nil
	case Tuple:
		return Value{}, nil
	default:
		return Value{}, fmt.Errorf("bruh you should not be here in decoding: unknown type error")
	}
}

// TODO create get... functions
// TODO if its get array/tuple function, pass index and take the element

func MakeUint8(value uint8) (Value, error) {
	bigInt := big.NewInt(int64(value))
	return MakeUint(bigInt, 8)
}

func MakeUint16(value uint16) (Value, error) {
	bigInt := big.NewInt(int64(value))
	return MakeUint(bigInt, 16)
}

func MakeUint32(value uint32) (Value, error) {
	bigInt := big.NewInt(int64(value))
	return MakeUint(bigInt, 32)
}

func MakeUint64(value uint64) (Value, error) {
	bigInt := big.NewInt(int64(0)).SetUint64(value)
	return MakeUint(bigInt, 64)
}

func MakeUint(value *big.Int, size uint16) (Value, error) {
	typeUint, err := MakeUintType(size)
	if err != nil {
		return Value{}, err
	}
	upperLimit := big.NewInt(0).Lsh(big.NewInt(1), uint(size))
	if value.Cmp(upperLimit) >= 0 {
		return Value{}, fmt.Errorf("passed value larger than uint size %d", size)
	}
	return Value{
		valueType: typeUint,
		value:     value,
	}, nil
}

func MakeUfixed(value *big.Rat, size uint16, precision uint16) (Value, error) {
	ufixedValueType, err := MakeUFixedType(size, precision)
	if err != nil {
		return Value{}, nil
	}
	denomSize := big.NewInt(0).Exp(
		big.NewInt(10), big.NewInt(int64(precision)),
		nil,
	)
	if value.Denom() != denomSize {
		return Value{}, fmt.Errorf("denominator size do not match")
	}
	numSize := big.NewInt(0).Lsh(big.NewInt(1), uint(size))
	if numSize.Cmp(value.Num()) <= 0 {
		return Value{}, fmt.Errorf("numerator size overflow")
	}
	return Value{
		valueType: ufixedValueType,
		value:     value,
	}, nil
}

func MakeString(value string) Value {
	return Value{
		valueType: MakeStringType(),
		value:     value,
	}
}

func MakeByte(value byte) Value {
	return Value{
		valueType: MakeByteType(),
		value:     value,
	}
}

func MakeAddress(value [32]byte) Value {
	return Value{
		valueType: MakeAddressType(),
		value:     value,
	}
}

func MakeDynamicArray(value []interface{}, elemType Type) Value {
	return Value{
		valueType: MakeDynamicArrayType(elemType),
		value:     value,
	}
}

func MakeStaticArray(value []interface{}, elemType Type) Value {
	return Value{
		valueType: MakeStaticArrayType(elemType, uint16(len(value))),
		value:     value,
	}
}

func MakeTuple(value []interface{}, tupleType []Type) (Value, error) {
	if len(value) != len(tupleType) {
		return Value{}, fmt.Errorf("tuple make: tuple element number unmatch with tuple type number")
	}
	if len(value) == 0 {
		return Value{}, fmt.Errorf("empty tuple")
	}
	return Value{
		valueType: MakeTupleType(tupleType),
		value:     value,
	}, nil
}

func MakeBool(value bool) Value {
	return Value{
		valueType: MakeBoolType(),
		value:     value,
	}
}

func GetUint8(value Value) (uint8, error) {
	if !(value.valueType.typeFromEnum == Uint && value.valueType.unsignedTypeSize <= 8) {
		return 0, fmt.Errorf("value type unmatch or size too large")
	}
	bigIntForm := value.value.(*big.Int)
	return uint8(bigIntForm.Uint64()), nil
}

func GetUint16(value Value) (uint16, error) {
	if !(value.valueType.typeFromEnum == Uint && value.valueType.unsignedTypeSize <= 16) {
		return 0, fmt.Errorf("value type unmatch or size too large")
	}
	bigIntForm := value.value.(*big.Int)
	return uint16(bigIntForm.Uint64()), nil
}

func GetUint32(value Value) (uint32, error) {
	if !(value.valueType.typeFromEnum == Uint && value.valueType.unsignedTypeSize <= 32) {
		return 0, fmt.Errorf("value type unmatch or size too large")
	}
	bigIntForm := value.value.(*big.Int)
	return uint32(bigIntForm.Uint64()), nil
}

func GetUint64(value Value) (uint64, error) {
	if !(value.valueType.typeFromEnum == Uint && value.valueType.unsignedTypeSize <= 64) {
		return 0, fmt.Errorf("value type unmatch or size too large")
	}
	bigIntForm := value.value.(*big.Int)
	return bigIntForm.Uint64(), nil
}

func GetUint(value Value) (*big.Int, error) {
	if value.valueType.typeFromEnum != Uint {
		return nil, fmt.Errorf("value type unmatch")
	}
	bigIntForm := value.value.(*big.Int)
	sizeThreshold := big.NewInt(0).Lsh(big.NewInt(1), uint(value.valueType.unsignedTypeSize))
	if sizeThreshold.Cmp(bigIntForm) <= 0 {
		return nil, fmt.Errorf("value is larger than uint size")
	}
	return bigIntForm, nil
}

func GetUfixed(value Value) (*big.Rat, error) {
	if value.valueType.typeFromEnum != Ufixed {
		return nil, fmt.Errorf("value type unmatch, should be ufixed")
	}
	ufixedForm := value.value.(*big.Rat)
	numinatorSize := big.NewInt(0).Lsh(big.NewInt(1), uint(value.valueType.unsignedTypeSize))
	denomSize := big.NewInt(0).Exp(
		big.NewInt(10), big.NewInt(int64(value.valueType.unsignedTypePrecision)),
		nil,
	)
	if denomSize.Cmp(ufixedForm.Denom()) != 0 || numinatorSize.Cmp(ufixedForm.Num()) <= 0 {
		return nil, fmt.Errorf("denominator size do not match, or numinator size overflow")
	}
	return ufixedForm, nil
}

func GetString(value Value) (string, error) {
	if value.valueType.typeFromEnum != String {
		return "", fmt.Errorf("value type unmatch, should be ufixed")
	}
	stringForm := value.value.(string)
	return stringForm, nil
}

func GetByte(value Value) (byte, error) {
	if value.valueType.typeFromEnum != Byte {
		return byte(0), fmt.Errorf("value type unmatch, should be bytes")
	}
	bytesForm := value.value.(byte)
	return bytesForm, nil
}

func GetAddress(value Value) ([32]byte, error) {
	if value.valueType.typeFromEnum != Address {
		return [32]byte{}, fmt.Errorf("value type unmatch, should be address")
	}
	addressForm := value.value.([32]byte)
	return addressForm, nil
}

func GetDynamicArrayByIndex(value Value, index uint16) (Value, error) {
	if value.valueType.typeFromEnum != ArrayDynamic {
		return Value{}, fmt.Errorf("value type unmatch, should be dynamic array")
	}
	elements := reflect.ValueOf(value.value)
	if int(index) >= elements.Len() {
		return Value{}, fmt.Errorf("dynamic array cannot get element: index out of scope")
	}
	return Value{
		valueType: value.valueType.childTypes[0],
		value:     elements.Index(int(index)).Interface(),
	}, nil
}

func GetStaticArrayByIndex(value Value, index uint16) (Value, error) {
	if value.valueType.typeFromEnum != ArrayStatic {
		return Value{}, fmt.Errorf("value type unmatch, should be static array")
	}
	if index >= value.valueType.staticLength {
		return Value{}, fmt.Errorf("static array cannot get element: index out of scope")
	}
	elements := reflect.ValueOf(value.value)
	return Value{
		valueType: value.valueType.childTypes[0],
		value:     elements.Index(int(index)).Interface(),
	}, nil
}

func GetTupleByIndex(value Value, index uint16) (Value, error) {
	if value.valueType.typeFromEnum != Tuple {
		return Value{}, fmt.Errorf("value type unmatch, should be tuple")
	}
	elements := reflect.ValueOf(value.value)
	if int(index) >= elements.Len() {
		return Value{}, fmt.Errorf("tuple cannot get element: index out of scope")
	}
	return Value{
		valueType: value.valueType.childTypes[index],
		value:     elements.Index(int(index)).Interface(),
	}, nil
}

func GetBool(value Value) (bool, error) {
	if value.valueType.typeFromEnum != Bool {
		return false, fmt.Errorf("value type unmatch, should be bool")
	}
	boolForm := value.value.(bool)
	return boolForm, nil
}
