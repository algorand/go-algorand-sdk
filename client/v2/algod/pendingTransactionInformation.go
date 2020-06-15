package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type PendingTransactionInformation struct {
	c    *Client
	txid string
	p    models.PendingTransactionInformationParams
}

func (s *PendingTransactionInformation) Max(max uint64) *PendingTransactionInformation {
	s.p.Max = max
	return s
}

// s.p.Format setter intentionally omitted: this SDK only uses msgpack.

func (s *PendingTransactionInformation) Do(ctx context.Context, headers ...*common.Header) (response models.PendingTransactionInfoResponse, stxn types.SignedTxn, err error) {
	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/transactions/pending/%s", s.txid), s.p, headers)
	stxn = response.Transaction
	return
}
