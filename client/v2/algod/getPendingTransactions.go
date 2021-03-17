package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type PendingTransactionsParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`

	// Max truncated number of transactions to display. If max=0, returns all pending
	// txns.
	Max uint64 `url:"max,omitempty"`
}

type PendingTransactions struct {
	c *Client

	p PendingTransactionsParams
}

// Max truncated number of transactions to display. If max=0, returns all pending
// txns.
func (s *PendingTransactions) Max(Max uint64) *PendingTransactions {
	s.p.Max = Max
	return s
}

func (s *PendingTransactions) Do(ctx context.Context, headers ...*common.Header) (total uint64, topTransactions []types.SignedTxn, err error) {
	s.p.Format = "msgpack"
	response := models.PendingTransactionsResponse{}
	err = s.c.getMsgpack(ctx, &response, "/v2/transactions/pending", s.p, headers)
	total = response.TotalTransactions
	topTransactions = response.TopTransactions
	return
}
