package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type ShutdownService struct {
	c *Client
	p models.ShutdownParams
}

func (s *ShutdownService) Timeout(timeout uint64) *ShutdownService {
	s.p.Timeout = timeout
	return s
}

func (s *ShutdownService) Do(ctx context.Context, headers ...*common.Header) error {
	return s.c.post(ctx, nil, "/shutdown", s.p, headers)
}
