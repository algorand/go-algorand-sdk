package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type Shutdown struct {
	c *Client
	p models.ShutdownParams
}

func (s *Shutdown) Timeout(timeout uint64) *Shutdown {
	s.p.Timeout = timeout
	return s
}

func (s *Shutdown) Do(ctx context.Context, headers ...*common.Header) error {
	return s.c.post(ctx, nil, "/shutdown", s.p, headers)
}
