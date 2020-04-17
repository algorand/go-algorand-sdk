package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type Status struct {
	c *Client
}

func (s *Status) Do(ctx context.Context, headers ...*common.Header) (status models.NodeStatus, err error) {
	err = s.c.get(ctx, &status, "/status", nil, headers)
	return
}
