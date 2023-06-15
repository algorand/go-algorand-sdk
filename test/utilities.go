package test

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	sdk_json "github.com/algorand/go-algorand-sdk/v2/encoding/json"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// VerifyResponse compares the actual response to the expected response.
func VerifyResponse(expectedFile string, actual string) error {
	jsonfile, err := os.Open(expectedFile)
	if err != nil {
		return err
	}
	fileBytes, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		return err
	}

	var expectedString string
	// No processing needed for json
	if strings.HasSuffix(expectedFile, ".json") {
		expectedString = string(fileBytes)
	}
	// Convert message pack to json for comparison
	if strings.HasSuffix(expectedFile, ".base64") {
		data, err := base64.StdEncoding.DecodeString(string(fileBytes))
		if err != nil {
			return fmt.Errorf("failed to decode '%s' from base64: %v", expectedFile, err)
		}
		generic := make(map[string]interface{})
		err = msgpack.Decode(data, generic)
		if err != nil {
			return fmt.Errorf("failed to decode '%s' from message pack: %v", expectedFile, err)
		}
		expectedString = string(sdk_json.EncodeStrict(generic))
	}

	if strings.HasSuffix(expectedFile, ".msgp") {
		generic := make(map[string]interface{})
		err = msgpack.Decode(fileBytes, generic)
		if err != nil {
			return fmt.Errorf("failed to decode '%s' from message pack: %v", expectedFile, err)
		}
		expectedString = string(sdk_json.EncodeStrict(generic))
	}

	err = EqualJSON2(expectedString, actual)
	if err != nil {
		fmt.Printf("EXPECTED:\n%v\n", expectedString)
		fmt.Printf("ACTUAL:\n%v\n", actual)
	}
	return err
}

// EqualJSON2 compares two json strings.
// returns true if considered equal, false otherwise.
// The error returns the difference.
// For reference: j1 is the baseline, j2 is the test
func EqualJSON2(j1, j2 string) (err error) {
	var expected map[string]interface{}
	err = json.Unmarshal([]byte(j1), &expected)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected: %w", err)
	}

	var actual map[string]interface{}
	err = json.Unmarshal([]byte(j2), &actual)
	if err != nil {
		return fmt.Errorf("failed to unmarshal actual: %w", err)
	}

	err = recursiveCompare("root", expected, actual)

	if err != nil {
		log.Printf("expected:\n%s", j1)
		log.Printf("actual:\n%s", j2)
	}
	return err
}

// ValueType is the type of the value as an enum.
type ValueType int

const (
	objectType ValueType = iota
	arrayType
	valueType
	numberType
	boolType
	stringType
	missingType
)

func getType(val interface{}) ValueType {
	if val == nil {
		return missingType
	}
	switch val.(type) {
	case map[string]interface{}:
		return objectType
	case []interface{}:
		return arrayType
	case string:
		return stringType
	case bool:
		return boolType
	case float64:
		return numberType
	case nil:
		return missingType
	default:
		return valueType
	}
}

// binaryOrStringEqual checks combinations of string / base64 decoded strings
// to see if the inputs are equal.
// The decoding process doesn't seem to distinguish between string and binary, but the encoding process
// does. So sometimes the string will be base64 encoded and need to compare against the decoded string
// value.
// There are some discrepancies in different algod / SDK types that causes
// encodings that need special handling:
// * Address is sometimes B32 encoded and sometimes B64 encoded.
// * BlockHash is sometimes B32 encoded (with a blk prefix) and sometimes B64 encoded.
func binaryOrStringEqual(s1, s2 string) bool {
	if s1 == s2 {
		return true
	}
	// S1 convert to S2
	{
		if val, err := base64.StdEncoding.DecodeString(s1); err == nil {
			if string(val) == s2 {
				return true
			}
			var addr types.Address
			if len(val) == len(addr[:]) {
				copy(addr[:], val)
				if addr.String() == s2 {
					return true
				}
			}

		}
		// parse blockhash
		var bh types.BlockHash
		if bh.UnmarshalText([]byte(s1)) == nil {
			if base64.StdEncoding.EncodeToString(bh[:]) == s2 {
				return true
			}
		}
	}

	// S2 convert to S1
	{
		if val, err := base64.StdEncoding.DecodeString(s2); err == nil {
			if string(val) == s1 {
				return true
			}
			var addr types.Address
			if len(val) == len(addr[:]) {
				copy(addr[:], val)
				if addr.String() == s1 {
					return true
				}
			}
		}
		var bh types.BlockHash
		if bh.UnmarshalText([]byte(s2)) == nil {
			if base64.StdEncoding.EncodeToString(bh[:]) == s1 {
				return true
			}
		}
	}
	return false
}

func sortArray(arr []interface{}, field string) {
	sort.SliceStable(arr, func(i, j int) bool {
		// literal type case
		if field == "" {
			vi := fmt.Sprintf("%v", arr[i])
			vj := fmt.Sprintf("%v", arr[j])
			return strings.Compare(vi, vj) < 0
		}
		// object case
		vali := arr[i].(map[string]interface{})[field]
		valj := arr[i].(map[string]interface{})[field]
		vi := fmt.Sprintf("%v", vali)
		vj := fmt.Sprintf("%v", valj)
		return strings.Compare(vi, vj) < 0
	})
}

func getFirstField(ob interface{}) string {
	if ob == nil || getType(ob) != objectType {
		return ""
	}
	for k := range ob.(map[string]interface{}) {
		return k
	}
	return ""
}

func recursiveCompare(field string, expected, actual interface{}) error {
	expectedType := getType(expected)
	actualType := getType(actual)

	// If both were nil, just return
	if expectedType == missingType && actualType == missingType {
		return nil
	}

	var keyType ValueType

	if expectedType == missingType || actualType == missingType {
		if expectedType == missingType {
			keyType = actualType
		}
		if actualType == missingType {
			keyType = expectedType
		}
	} else {
		// If both are present, make sure they are the same
		if expectedType != actualType {
			return errors.New("Type mismatch")
		}
		keyType = expectedType
	}

	switch keyType {
	case arrayType:
		var expectedArr []interface{}
		var actualArr []interface{}

		expectedSize := 0
		if expectedType != missingType {
			expectedArr = expected.([]interface{})
			expectedSize = len(expectedArr)
		}
		actualSize := 0
		if actualType != missingType {
			actualArr = actual.([]interface{})
			actualSize = len(actualArr)
		}

		if expectedSize != actualSize {
			return fmt.Errorf("failed to match array sizes: %s", field)
		}

		//sortField := getFirstField(expected)
		//sortArray(actualArr, sortField)
		//sortArray(expectedArr, sortField)

		// n^2 baby! Just make sure every element has a match somewhere in the other list.
		// getting the sort function to work would be better...
		var err error
		for i := 0; i < expectedSize; i++ {
			for j := 0; j < len(actualArr); j++ {
				err = recursiveCompare(fmt.Sprintf("%s[%d]", field, i), expectedArr[i], actualArr[j])
				if err == nil {
					break
				}
			}
			if err != nil {
				return err
			}
		}
		return err

	case objectType:
		//log.Printf("%s{...} - object\n", field)

		// Recursively compare each key value
		// Pass nil's to the compare function to handle zero values on a type by type basis.

		// Go happily creates complex zero value objects, so go ahead and recursively compare nil against defaults

		// If they are both missing what are we even doing here. Return with no error.
		if expectedType == missingType && actualType == missingType {
			return nil
		}

		var expectedObject map[string]interface{}
		var actualObject map[string]interface{}

		keys := make(map[string]bool)
		if expectedType != missingType {
			expectedObject = expected.(map[string]interface{})
			for k := range expectedObject {
				keys[k] = true
			}
		}
		if actualType != missingType {
			actualObject = actual.(map[string]interface{})
			for k := range actualObject {
				keys[k] = true
			}
		}
		for k := range keys {
			err := recursiveCompare(fmt.Sprintf("%s.%s", field, k), expectedObject[k], actualObject[k])
			if err != nil {
				return err
			}
		}

	case numberType:
		// Compare numbers, if missing treat as zero
		expectedNum := float64(0)
		if expectedType != missingType {
			expectedNum = expected.(float64)
		}
		actualNum := float64(0)
		if actualType != missingType {
			actualNum = actual.(float64)
		}
		//log.Printf("%s - number %f == %f\n", field, expectedNum, actualNum)
		if expectedNum != actualNum {
			return fmt.Errorf("failed to match field %s, %f != %f", field, expectedNum, actualNum)
		}

	case boolType:
		// Compare bools, if missing treat as false
		expectedBool := false
		if expectedType != missingType {
			expectedBool = expected.(bool)
		}
		actualBool := false
		if actualType != missingType {
			actualBool = actual.(bool)
		}
		//log.Printf("%s - bool %t == %t\n", field, expectedBool, actualBool)
		if expectedBool != actualBool {
			return fmt.Errorf("failed to match field %s, %t != %t", field, expectedBool, actualBool)
		}

	case stringType:
		// Compare strings, if missing treat as an empty string.
		// Note: I think binary ends up in here, it may need some special handling.
		expectedStr := ""
		if expectedType != missingType {
			expectedStr = expected.(string)
		}
		actualStr := ""
		if actualType != missingType {
			actualStr = actual.(string)
		}

		//log.Printf("%s - string %s == %s\n", field, expectedStr, actualStr)
		if !binaryOrStringEqual(expectedStr, actualStr) {
			return fmt.Errorf("failed to match field %s, %s != %s", field, expectedStr, actualStr)
		}

	default:
		return fmt.Errorf("unhandled type %v at %s", keyType, field)
	}

	return nil
}

func sliceOfBytesEqual(expected [][]byte, actual [][]byte) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("expected length (%d) does not match actual length (%d)", len(expected), len(actual))
	}

	for i, expectedElement := range expected {
		actualElement := actual[i]
		if !bytes.Equal(expectedElement, actualElement) {
			return fmt.Errorf("elements at index %d are unequal. Expected %s, got %s", i, hex.EncodeToString(expectedElement), hex.EncodeToString(actualElement))
		}
	}

	return nil
}
