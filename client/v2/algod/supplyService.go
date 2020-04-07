package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type SupplyService struct {
	c *Client
}

func (s *SupplyService) Do(ctx context.Context, headers ...*common.Header) (supply models.Supply, err error) {
	err = s.c.get(ctx, &supply, "/ledger/supply", nil, headers)
	return
}
