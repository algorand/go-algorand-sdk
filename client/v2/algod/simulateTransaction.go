package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// SimulateTransactionParams contains all of the query parameters for url serialization.
type SimulateTransactionParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded. If
	// not provided, defaults to JSON.
	Format string `url:"format,omitempty"`
}

// SimulateTransaction simulates a raw transaction or transaction group as it would
// be evaluated on the network. The simulation will use blockchain state from the
// latest committed round.
type SimulateTransaction struct {
	c *Client

	request models.SimulateRequest

	p SimulateTransactionParams
}

// Do performs the HTTP request
func (s *SimulateTransaction) Do(ctx context.Context, headers ...*common.Header) (response models.SimulateResponse, err error) {
	err = s.c.post(ctx, &response, "/v2/transactions/simulate", s.p, headers, s.request)
	return
}
