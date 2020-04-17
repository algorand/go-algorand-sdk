package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type PendingTransactions struct {
	c *Client
}

func (s *PendingTransactions) Do(ctx context.Context, headers ...*common.Header) (total uint64, topTransactions []types.SignedTxn, err error) {
	response := models.PendingTransactionsResponse{}
	err = s.c.get(ctx, &response, "/transactions/pending", nil, headers)
	total = response.TotalTransactions
	for _, b64SignedTxn := range response.TopTransactions {
		var signedTxn types.SignedTxn
		err = signedTxn.FromBase64String(b64SignedTxn)
		if err != nil {
			return
		}
		topTransactions = append(topTransactions, signedTxn)
	}
	return
}
