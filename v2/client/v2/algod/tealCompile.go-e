package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// TealCompile given TEAL source code in plain text, return base64 encoded program
// bytes and base32 SHA512_256 hash of program bytes (Address style). This endpoint
// is only enabled when a node's configureation file sets EnableDeveloperAPI to
// true.
type TealCompile struct {
	c *Client

	source []byte
}

// Do performs the HTTP request
func (s *TealCompile) Do(ctx context.Context, headers ...*common.Header) (response models.CompileResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/teal/compile", s.source, headers)
	return
}
