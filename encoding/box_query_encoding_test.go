package encoding

import (
	"encoding/json"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stretchr/testify/require"
)

type example struct {
	source           []byte
	expectedEncoding string
}

var e = example{[]byte("tkÿÿ"), "b64:dGvDv8O/"}

func TestEncode(t *testing.T) {
	// Encode bytes
	require.Equal(t,
		e.expectedEncoding,
		EncodeBytesForBoxQuery(e.source),
	)

	// Encode Box
	b := models.Box{Name: e.source}
	actual := EncodeBoxForBoxQuery(b)
	require.Equal(t,
		e.expectedEncoding,
		actual,
	)

	// Encode BoxDescriptor
	bd := models.BoxDescriptor{Name: e.source}
	actual = EncodeBoxDescriptorForBoxQuery(bd)
	require.Equal(t,
		e.expectedEncoding,
		actual,
	)

	// Encode BoxReference
	require.Equal(t,
		e.expectedEncoding,
		EncodeBoxReferenceForBoxQuery(types.BoxReference{Name: e.source}))

	// Encode AppBoxReference
	require.Equal(t,
		e.expectedEncoding,
		EncodeAppBoxReferenceForBoxQuery(types.AppBoxReference{Name: e.source}))
}

func TestEncodeFromJSON(t *testing.T) {
	jsonEncoding := []byte(`{"name":"ZXhhbXBsZSBib3ggbmFtZQ==","value":"dGvDv8O/"}`)

	var box models.Box
	err := json.Unmarshal(jsonEncoding, &box)
	require.NoError(t, err)

	actualEncodedName := EncodeBoxForBoxQuery(box)
	expectedEncodedName := "b64:ZXhhbXBsZSBib3ggbmFtZQ=="
	require.Equal(t, expectedEncodedName, actualEncodedName)
}
