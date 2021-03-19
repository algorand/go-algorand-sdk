package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

// HealthCheck returns OK if healthy.
type HealthCheck struct {
	c *Client
}

// Do performs the HTTP request
func (s *HealthCheck) Do(ctx context.Context, headers ...*common.Header) error {
	return s.c.get(ctx, nil, "/health", nil, headers)
}
