package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

// SetSyncRound sets the minimum sync round on the ledger.
type SetSyncRound struct {
	c *Client

	round uint64
}

// Do performs the HTTP request
func (s *SetSyncRound) Do(ctx context.Context, headers ...*common.Header) (response string, err error) {
	err = s.c.post(ctx, &response, fmt.Sprintf("/v2/ledger/sync/%s", common.EscapeParams(s.round)...), nil, headers, nil)
	return
}
