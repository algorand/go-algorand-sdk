package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// TealDisassemble given the program bytes, return the TEAL source code in plain
// text. This endpoint is only enabled when a node's configuration file sets
// EnableDeveloperAPI to true.
type TealDisassemble struct {
	c *Client

	source []byte
}

// Do performs the HTTP request
func (s *TealDisassemble) Do(ctx context.Context, headers ...*common.Header) (response models.DisassembleResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/teal/disassemble", nil, headers, s.source)
	return
}
