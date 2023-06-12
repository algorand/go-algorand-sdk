package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// GetLedgerStateDeltaParams contains all of the query parameters for url serialization.
type GetLedgerStateDeltaParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded. If
	// not provided, defaults to JSON.
	Format string `url:"format,omitempty"`
}

// GetLedgerStateDelta get ledger deltas for a round.
type GetLedgerStateDelta struct {
	c *Client

	round uint64

	p GetLedgerStateDeltaParams
}

// Do performs the HTTP request
func (s *GetLedgerStateDelta) Do(ctx context.Context, headers ...*common.Header) (response types.LedgerStateDelta, err error) {
	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/deltas/%s", common.EscapeParams(s.round)...), s.p, headers)
	return
}
