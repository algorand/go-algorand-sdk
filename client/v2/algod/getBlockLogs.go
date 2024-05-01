package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetBlockLogs get all of the logs from outer and inner app calls in the given
// round
type GetBlockLogs struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *GetBlockLogs) Do(ctx context.Context, headers ...*common.Header) (response models.BlockLogsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/blocks/%s/logs", common.EscapeParams(s.round)...), nil, headers)
	return
}
