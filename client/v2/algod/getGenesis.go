package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

type GetGenesis struct {
	c *Client
}

func (s *GetGenesis) Do(ctx context.Context, headers ...*common.Header) (response string, err error) {
	err = s.c.get(ctx, &response, "/genesis", nil, headers)
	return
}
