package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

// PendingTransactionsByAddressParams contains all of the query parameters for url serialization.
type PendingTransactionsByAddressParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`

	// Max truncated number of transactions to display. If max=0, returns all pending
	// txns.
	Max uint64 `url:"max,omitempty"`
}

// PendingTransactionsByAddress get the list of pending transactions by address,
// sorted by priority, in decreasing order, truncated at the end at MAX. If MAX =
// 0, returns all pending transactions.
type PendingTransactionsByAddress struct {
	c *Client

	address string

	p PendingTransactionsByAddressParams
}

// Max truncated number of transactions to display. If max=0, returns all pending
// txns.
func (s *PendingTransactionsByAddress) Max(Max uint64) *PendingTransactionsByAddress {
	s.p.Max = Max

	return s
}

// Do performs the HTTP request
func (s *PendingTransactionsByAddress) Do(ctx context.Context, headers ...*common.Header) (total uint64, topTransactions []types.SignedTxn, err error) {
	s.p.Format = "msgpack"
	response := models.PendingTransactionsResponse{}
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/accounts/%s/transactions/pending", s.address), s.p, headers)
	total = response.TotalTransactions
	topTransactions = response.TopTransactions
	return
}
