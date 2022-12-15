package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetSyncRound gets the minimum sync round for the ledger.
type GetSyncRound struct {
	c *Client
}

// Do performs the HTTP request
func (s *GetSyncRound) Do(ctx context.Context, headers ...*common.Header) (response models.GetSyncRoundResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/ledger/sync", nil, headers)
	return
}
