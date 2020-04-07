package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type StatusAfterBlockService struct {
	c     *Client
	round uint64
}

func (s *StatusAfterBlockService) Do(ctx context.Context, headers ...*common.Header) (status models.NodeStatus, err error) {
	err = s.c.get(ctx, &status, fmt.Sprintf("/status/wait-for-block-after/%d", s.round), nil, headers)
	return
}
