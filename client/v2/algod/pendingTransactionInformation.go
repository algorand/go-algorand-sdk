package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type PendingTransactionInformationParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

type PendingTransactionInformation struct {
	c *Client

	txid string

	p PendingTransactionInformationParams
}

func (s *PendingTransactionInformation) Do(ctx context.Context, headers ...*common.Header) (response models.PendingTransactionInfoResponse, stxn types.SignedTxn, err error) {
	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/transactions/pending/%s", s.txid), s.p, headers)
	stxn = response.Transaction
	return
}
