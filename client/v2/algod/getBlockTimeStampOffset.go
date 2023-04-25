package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetBlockTimeStampOffset gets the current timestamp offset.
type GetBlockTimeStampOffset struct {
	c *Client
}

// Do performs the HTTP request
func (s *GetBlockTimeStampOffset) Do(ctx context.Context, headers ...*common.Header) (response models.GetBlockTimeStampOffsetResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/devmode/blocks/offset", nil, headers)
	return
}
