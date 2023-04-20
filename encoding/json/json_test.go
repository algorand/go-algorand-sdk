package json

import (
	"bytes"
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
		require.NoError(t, err)
		assert.Equal(t, obj, decoded)
	})

	t.Run("strict decode, pass", func(t *testing.T) {
		// strict decode test
		decoder := NewDecoder(bytes.NewReader(encodedOb))
		var decoded object
		err := decoder.Decode(&decoded)
		require.NoError(t, err)
		assert.Equal(t, obj, decoded)
	})

	t.Run("strict decode subset, fail", func(t *testing.T) {
		// strict decode test
		decoder := NewDecoder(bytes.NewReader(encodedOb))
		var decoded subsetObject
		err := decoder.Decode(&decoded)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no matching struct field found when decoding stream map with key name")
	})

	t.Run("lenient decode subset, pass", func(t *testing.T) {
		// strict decode test
		decoder := NewLenientDecoder(bytes.NewReader(encodedOb))
		var decoded subsetObject
		err := decoder.Decode(&decoded)
		require.NoError(t, err)
		assert.Equal(t, obj.subsetObject, decoded)
	})

	t.Run("original encode map key as string", func(t *testing.T) {
		intMap := map[int]string{
			0: "int key",
		}
		data := string(Encode(intMap))
		assert.NotContains(t, data, "\"0\":")
	})

	t.Run("strict encode map key as string", func(t *testing.T) {
		intMap := map[int]string{
			0: "int key",
		}
		data := string(EncodeStrict(intMap))
		assert.NotContains(t, data, "0:")
	})

	t.Run("strict encode map interface key as string", func(t *testing.T) {
		t.Skip("There is a bug in go-codec with MapKeyAsString = true and Canonical = true")
		intMap := map[interface{}]interface{}{
			0: "int key",
		}
		data := string(EncodeStrict(intMap))
		assert.NotContains(t, data, "0:")
	})
}
