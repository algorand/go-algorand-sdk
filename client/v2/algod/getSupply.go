package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// Supply get the current supply reported by the ledger.
type Supply struct {
	c *Client
}

// Do performs the HTTP request
func (s *Supply) Do(ctx context.Context, headers ...*common.Header) (response models.Supply, err error) {
	err = s.c.get(ctx, &response, "/v2/ledger/supply", nil, headers)
	return
}
