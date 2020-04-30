package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type PendingTransactions struct {
	c *Client
	p models.PendingTransactionInformationParams
}

func (s *PendingTransactions) Max(max uint64) *PendingTransactions {
	s.p.Max = max
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
