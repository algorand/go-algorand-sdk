package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupTransactions struct {
	c *Client

	txid string
}

func (s *LookupTransactions) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/transactions/%v", s.txid), nil, headers)
	return
}
