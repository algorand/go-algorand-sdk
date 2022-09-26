package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetBlockHash get the block hash for the block on the given round.
type GetBlockHash struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *GetBlockHash) Do(ctx context.Context, headers ...*common.Header) (response models.BlockHashResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%s/hash", common.EscapeParams(s.round)...), nil, headers)
	return
}
