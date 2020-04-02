package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type PendingTransactionInformationService struct {
	c    *Client
	txid string
	p    models.PendingTransactionInformationParams
}

func (s *PendingTransactionInformationService) Max(max uint64) *PendingTransactionInformationService {
	s.p.Max = max
	return s
}

// s.p.Format setter intentionally omitted: this SDK only uses msgpack.

func (s *PendingTransactionInformationService) Do(ctx context.Context, headers ...*common.Header) (response models.PendingTransactionInfoResponse, stxn types.SignedTxn, err error) {
	s.p.Format = "msgpack"
	err = s.c.get(ctx, &response, fmt.Sprintf("/transactions/pending/%s", s.txid), s.p, headers)
	if err != nil {
		return
	}
	err = stxn.FromBase64String(response.PendingTransactionBase64)
	return
}
