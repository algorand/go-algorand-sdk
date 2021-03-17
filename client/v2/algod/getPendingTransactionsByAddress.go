package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type PendingTransactionInformationByAddressParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`

	// Max truncated number of transactions to display. If max=0, returns all pending
	// txns.
	Max uint64 `url:"max,omitempty"`
}

type PendingTransactionInformationByAddress struct {
	c *Client

	address string

	p PendingTransactionInformationByAddressParams
}

// Max truncated number of transactions to display. If max=0, returns all pending
// txns.
func (s *PendingTransactionInformationByAddress) Max(Max uint64) *PendingTransactionInformationByAddress {
	s.p.Max = Max
	return s
}

func (s *PendingTransactionInformationByAddress) Do(ctx context.Context, headers ...*common.Header) (total uint64, topTransactions []types.SignedTxn, err error) {
	s.p.Format = "msgpack"
	response := models.PendingTransactionsResponse{}
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/accounts/%s/transactions/pending", s.address), s.p, headers)
	total = response.TotalTransactions
	topTransactions = response.TopTransactions
	return
}
