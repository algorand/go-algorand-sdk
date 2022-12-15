package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetLedgerStateDelta get ledger deltas for a round.
type GetLedgerStateDelta struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *GetLedgerStateDelta) Do(ctx context.Context, headers ...*common.Header) (response models.LedgerStateDelta, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/deltas/%s", common.EscapeParams(s.round)...), nil, headers)
	return
}
