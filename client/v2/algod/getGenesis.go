package algod

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

// GetGenesis returns the entire genesis file in json.
type GetGenesis struct {
	c *Client
}

// Do performs the HTTP request
func (s *GetGenesis) Do(ctx context.Context, headers ...*common.Header) (response string, err error) {
	err = s.c.get(ctx, &response, "/genesis", nil, headers)
	return
}
