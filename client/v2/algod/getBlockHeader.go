package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetBlockHeaderParams contains all of the query parameters for url serialization.
type GetBlockHeaderParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded. If
	// not provided, defaults to JSON.
	Format string `url:"format,omitempty"`
}

// GetBlockHeader get the block header for the block on the given round.
type GetBlockHeader struct {
	c *Client

	round uint64

	p GetBlockHeaderParams
}

// Do performs the HTTP request
func (s *GetBlockHeader) Do(ctx context.Context, headers ...*common.Header) (response models.BlockHeaderResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%s/header", common.EscapeParams(s.round)...), s.p, headers)
	return
}
