package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// GetLedgerStateDeltaForTransactionGroupParams contains all of the query parameters for url serialization.
type GetLedgerStateDeltaForTransactionGroupParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded. If
	// not provided, defaults to JSON.
	Format string `url:"format,omitempty"`
}

// GetLedgerStateDeltaForTransactionGroup get a ledger delta for a given
// transaction group.
type GetLedgerStateDeltaForTransactionGroup struct {
	c *Client

	id string

	p GetLedgerStateDeltaForTransactionGroupParams
}

// Do performs the HTTP request
func (s *GetLedgerStateDeltaForTransactionGroup) Do(ctx context.Context, headers ...*common.Header) (response types.LedgerStateDelta, err error) {
	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/deltas/txn/group/%s", common.EscapeParams(s.id)...), s.p, headers)
	return
}
