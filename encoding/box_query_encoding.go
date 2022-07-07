package encoding

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

const encodingBase64Prefix = "b64:"

// EncodeBytesForBoxQuery provides a convenience method to string encode box names for use with Box search APIs (e.g. GetApplicationBoxByName).
func EncodeBytesForBoxQuery(xs []byte) string {
	return encodingBase64Prefix + base64.StdEncoding.EncodeToString(xs)
}

// EncodeBoxForBoxQuery provides a convenience method to string encode box names for use with Box search APIs (e.g. GetApplicationBoxByName).
func EncodeBoxForBoxQuery(b models.Box) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(b.Name))
	if err != nil {
		return "", err
	}

	return EncodeBytesForBoxQuery(decoded), nil
}

// EncodeBoxReferenceForBoxQuery provides a convenience method to string encode box names for use with Box search APIs (e.g. GetApplicationBoxByName).
func EncodeBoxReferenceForBoxQuery(br types.BoxReference) string {
	return EncodeBytesForBoxQuery(br.Name)
}

func EncodeAppBoxReferenceForBoxQuery(abr types.AppBoxReference) string {
	return EncodeBytesForBoxQuery(abr.Name)
}
