package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
)

// GetReady returns OK if healthy and fully caught up.
type GetReady struct {
	c *Client
}

// Do performs the HTTP request
func (s *GetReady) Do(ctx context.Context, headers ...*common.Header) (response string, err error) {
	err = s.c.get(ctx, &response, "/ready", nil, headers)
	return
}
