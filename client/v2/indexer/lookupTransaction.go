package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupTransaction lookup a single transaction.
type LookupTransaction struct {
	c *Client

	txid string
}

// Do performs the HTTP request
func (s *LookupTransaction) Do(ctx context.Context, headers ...*common.Header) (response models.TransactionResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/transactions/%s", common.EscapeParams(s.txid)...), nil, headers)
	return
}
