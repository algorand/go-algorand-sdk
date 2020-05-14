package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type HealthCheck struct {
	c *Client
}

func (s *HealthCheck) Do(ctx context.Context, headers ...*common.Header) (healthCheck models.HealthCheckResponse, err error) {
	var response models.HealthCheckResponse
	err = s.c.get(ctx, &response, "/health", nil, headers)
	healthCheck = models.HealthCheckResponse(response)
	return
}
