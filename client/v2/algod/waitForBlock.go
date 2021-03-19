package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type StatusAfterBlock struct {
	c *Client

	round uint64
}

func (s *StatusAfterBlock) Do(ctx context.Context, headers ...*common.Header) (response models.NodeStatus, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/status/wait-for-block-after/%v", s.round), nil, headers)
	return
}
