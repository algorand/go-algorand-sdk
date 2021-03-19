package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

type HealthCheck struct {
	c *Client
}

func (s *HealthCheck) Do(ctx context.Context, headers ...*common.Header) error {
	return s.c.get(ctx, nil, "/health", nil, headers)
}
