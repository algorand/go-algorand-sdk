package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// StatusAfterBlock waits for a block to appear after round {round} and returns the
// node's status at the time. There is a 1 minute timeout, when reached the current
// status is returned regardless of whether or not it is the round after the given
// round.
type StatusAfterBlock struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *StatusAfterBlock) Do(ctx context.Context, headers ...*common.Header) (response models.NodeStatus, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/status/wait-for-block-after/%s", common.EscapeParams(s.round)...), nil, headers)
	return
}
