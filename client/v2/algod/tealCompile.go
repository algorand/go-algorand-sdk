package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// TealCompileParams contains all of the query parameters for url serialization.
type TealCompileParams struct {

	// Sourcemap when set to `true`, returns the source map of the program as a JSON.
	// Defaults to `false`.
	Sourcemap bool `url:"sourcemap,omitempty"`
}

// TealCompile given TEAL source code in plain text, return base64 encoded program
// bytes and base32 SHA512_256 hash of program bytes (Address style). This endpoint
// is only enabled when a node's configuration file sets EnableDeveloperAPI to
// true.
type TealCompile struct {
	c *Client

	source []byte

	p TealCompileParams
}

// Sourcemap when set to `true`, returns the source map of the program as a JSON.
// Defaults to `false`.
func (s *TealCompile) Sourcemap(Sourcemap bool) *TealCompile {
	s.p.Sourcemap = Sourcemap

	return s
}

// Do performs the HTTP request
func (s *TealCompile) Do(ctx context.Context, headers ...*common.Header) (response models.CompileResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/teal/compile", s.p, headers, s.source)
	return
}
