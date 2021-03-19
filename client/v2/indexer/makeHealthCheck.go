package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// MakeHealthCheck returns 200 if healthy.
type MakeHealthCheck struct {
	c *Client
}

// Do performs the HTTP request
func (s *MakeHealthCheck) Do(ctx context.Context, headers ...*common.Header) (response models.HealthCheck, err error) {
	err = s.c.get(ctx, &response, "/health", nil, headers)
	return
}
