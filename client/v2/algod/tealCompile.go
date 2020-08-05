package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// TealCompile /v2/teal/compile
// Given TEAL source code in plain text, return base64 encoded program bytes and
// base32 SHA512_256 hash of program bytes (Address style).
type TealCompile struct {
	c      *Client
	source []byte
}

// Do performs HTTP request
func (s *TealCompile) Do(ctx context.Context,
	headers ...*common.Header) (response models.CompileResponse, err error) {
	err = s.c.post(ctx, &response,
		"/v2/teal/compile", s.source, headers)
	return
}
