package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type PendingTransactionInformationByAddress struct {
	c       *Client
	address string
	p       models.GetPendingTransactionsByAddressParams
}

func (s *PendingTransactionInformationByAddress) Max(max uint64) *PendingTransactionInformationByAddress {
	s.p.Max = max
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
