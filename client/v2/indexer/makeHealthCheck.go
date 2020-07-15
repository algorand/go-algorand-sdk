package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /health
 *
 */
type MakeHealthCheck struct {
	c *Client
}

func (s *MakeHealthCheck) Do(ctx context.Context,
	headers ...*common.Header) (response models.HealthCheckResponse, err error) {
	err = s.c.get(ctx, &response,
		"/health", nil, headers)
	return
}
