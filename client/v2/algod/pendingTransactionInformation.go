package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

// PendingTransactionInformationParams contains all of the query parameters for url serialization.
type PendingTransactionInformationParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

// PendingTransactionInformation given a transaction ID of a recently submitted
// transaction, it returns information about it. There are several cases when this
// might succeed:
// - transaction committed (committed round > 0)
// - transaction still in the pool (committed round = 0, pool error = "")
// - transaction removed from pool due to error (committed round = 0, pool error !=
// "")
// Or the transaction may have happened sufficiently long ago that the node no
// longer remembers it, and this will return an error.
type PendingTransactionInformation struct {
	c *Client

	txid string

	p PendingTransactionInformationParams
}

// Do performs the HTTP request
func (s *PendingTransactionInformation) Do(ctx context.Context, headers ...*common.Header) (response models.PendingTransactionInfoResponse, stxn types.SignedTxn, err error) {
	s.p.Format = "msgpack"
	err = s.c.getMsgpack(ctx, &response, fmt.Sprintf("/v2/transactions/pending/%s", s.txid), s.p, headers)
	stxn = response.Transaction
	return
}
