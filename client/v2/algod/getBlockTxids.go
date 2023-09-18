package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetBlockTxids get the top level transaction IDs for the block on the given
// round.
type GetBlockTxids struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *GetBlockTxids) Do(ctx context.Context, headers ...*common.Header) (response models.BlockTxidsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%s/txids", common.EscapeParams(s.round)...), nil, headers)
	return
}
