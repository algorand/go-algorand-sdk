package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// HealthCheck returns 200 if healthy.
type HealthCheck struct {
	c *Client
}

// Do performs the HTTP request
func (s *HealthCheck) Do(ctx context.Context, headers ...*common.Header) (response models.HealthCheckResponse, err error) {
	err = s.c.get(ctx, &response, "/health", nil, headers)
	return
}
