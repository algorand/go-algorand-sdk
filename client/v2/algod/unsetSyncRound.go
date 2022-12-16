package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
)

// UnsetSyncRound unset the ledger sync round.
type UnsetSyncRound struct {
	c *Client
}

// Do performs the HTTP request
func (s *UnsetSyncRound) Do(ctx context.Context, headers ...*common.Header) (response string, err error) {
	err = s.c.delete(ctx, &response, "/v2/ledger/sync", nil, headers)
	return
}
