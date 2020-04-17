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
	response := models.PendingTransactionsResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("/accounts/%s/transactions/pending", s.address), s.p, headers)
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
