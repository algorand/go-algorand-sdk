package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

type HealthCheckService struct {
	c *Client
}

func (s *HealthCheckService) Do(ctx context.Context, headers ...*common.Header) error {
	return s.c.get(ctx, nil, "/health", nil, headers)
}
