package encoding

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
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
	b := models.Box{Name: []byte(base64.StdEncoding.EncodeToString(e.source))}
	actual, err := EncodeBoxForBoxQuery(b)
	require.NoError(t, err)
	require.Equal(t,
		e.expectedEncoding,
		actual,
	)

	// Encode BoxDescriptor
	bd := models.BoxDescriptor{Name: []byte(base64.StdEncoding.EncodeToString(e.source))}
	actual, err = EncodeBoxDescriptorForBoxQuery(bd)
	require.NoError(t, err)
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
