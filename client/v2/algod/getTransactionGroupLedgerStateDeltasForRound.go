package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// GetTransactionGroupLedgerStateDeltasForRoundParams contains all of the query parameters for url serialization.
type GetTransactionGroupLedgerStateDeltasForRoundParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded. If
	// not provided, defaults to JSON.
	Format string `url:"format,omitempty"`
}

// GetTransactionGroupLedgerStateDeltasForRound get ledger deltas for transaction
// groups in a given round.
type GetTransactionGroupLedgerStateDeltasForRound struct {
	c *Client

	round uint64

	p GetTransactionGroupLedgerStateDeltasForRoundParams
}

// Do performs the HTTP request
func (s *GetTransactionGroupLedgerStateDeltasForRound) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionGroupLedgerStateDeltasForRoundResponse, err error) {
	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/deltas/%s/txn/group", common.EscapeParams(s.round)...), s.p, headers)
	return
}
