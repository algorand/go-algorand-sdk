package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type Versions struct {
	c *Client
}

func (s *Versions) Do(ctx context.Context, headers ...*common.Header) (response models.Version, err error) {
	err = s.c.get(ctx, &response, "/versions", nil, headers)
	return
}
