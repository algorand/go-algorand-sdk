package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type object struct {
	subsetObject
	Name string `codec:"name"`
}

type subsetObject struct {
	Data string `codec:"data"`
}

var obj object
var encodedOb []byte

func init() {
	obj = object{
		subsetObject: subsetObject{Data: "data"},
		Name:         "name",
	}
	encodedOb = Encode(obj)
}

func TestBasicEncodeDecode(t *testing.T) {
	// basic encode/decode test.
	var decoded object
	err := Decode(encodedOb, &decoded)
	require.NoError(t, err)
	assert.Equal(t, obj, decoded)
}

func TestDecode(t *testing.T) {
	decoder := NewDecoder(bytes.NewReader(encodedOb))
	var decoded object
	err := decoder.Decode(&decoded)
	require.NoError(t, err)
	assert.Equal(t, obj, decoded)
}

func TestSubsetDecode(t *testing.T) {
	decoder := NewDecoder(bytes.NewReader(encodedOb))
	var decoded subsetObject
	err := decoder.Decode(&decoded)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no matching struct field found when decoding stream map with key name")
}

func TestLenientDecode(t *testing.T) {
	decoder := NewLenientDecoder(bytes.NewReader(encodedOb))
	var decoded subsetObject
	err := decoder.Decode(&decoded)
	require.NoError(t, err)
	assert.Equal(t, obj.subsetObject, decoded)
}

func TestEncodeMaapKeyAsString(t *testing.T) {
	intMap := map[int]string{
		0: "int key",
	}
	data := string(Encode(intMap))
	assert.NotContains(t, data, `"0"`)
}

func TestStrictEncodeMapIntKeyAsString(t *testing.T) {
	intMap := map[int]string{
		0: "int key",
	}
	data := string(EncodeStrict(intMap))
	assert.NotContains(t, data, "0:")
}

func TestStrictEncodeMapInterfaceKeyAsString(t *testing.T) {
	intMap := map[interface{}]interface{}{
		0: "int key",
	}
	data := string(EncodeStrict(intMap))
	assert.Contains(t, data, `"0"`)
}

func TestStructKeyEncode(t *testing.T) {
	type KeyStruct struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}
	type TestStruct struct {
		Complex map[KeyStruct]string `json:"complex"`
	}

	data := TestStruct{
		Complex: map[KeyStruct]string{
			KeyStruct{
				Key1: "key1",
				Key2: "key2",
			}: "value",
		},
	}

	encoded := Encode(data)
	fmt.Println(string(encoded))

	var data2 TestStruct
	err := Decode(encoded, &data2)
	assert.NoError(t, err)
	assert.Equal(t, data, data2)

	// Unfortunately, still an error
	var data3 TestStruct
	err = json.NewDecoder(bytes.NewReader(encoded)).Decode(&data3)
	assert.Error(t, err)
}
