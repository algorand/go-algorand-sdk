package msgpack

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type object struct {
	subsetObject
	Name string `codec:"name"`
}

type subsetObject struct {
	Data string `codec:"data"`
}

func TestDecode(t *testing.T) {
	obj := object{
		subsetObject: subsetObject{Data: "data"},
		Name:         "name",
	}
	encodedOb := Encode(obj)

	t.Run("basic encode/decode test", func(t *testing.T) {
		// basic encode/decode test.
		var decoded object
		err := Decode(encodedOb, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, obj, decoded)
	})

	t.Run("strict decode, pass", func(t *testing.T) {
		// strict decode test
		decoder := NewDecoder(bytes.NewReader(encodedOb))
		var decoded object
		err := decoder.Decode(&decoded)
		assert.NoError(t, err)
		assert.Equal(t, obj, decoded)
	})

	t.Run("strict decode subset, fail", func(t *testing.T) {
		// strict decode test
		decoder := NewDecoder(bytes.NewReader(encodedOb))
		var decoded subsetObject
		err := decoder.Decode(&decoded)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no matching struct field found when decoding stream map with key name")
	})

	t.Run("lenient decode subset, pass", func(t *testing.T) {
		// strict decode test
		decoder := NewLenientDecoder(bytes.NewReader(encodedOb))
		var decoded subsetObject
		err := decoder.Decode(&decoded)
		assert.NoError(t, err)
		assert.Equal(t, obj.subsetObject, decoded)
	})
}
