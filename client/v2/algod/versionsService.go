package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type VersionsService struct {
	c *Client
}

func (s *VersionsService) Do(ctx context.Context, headers ...*common.Header) (response models.Version, err error) {
	err = s.c.get(ctx, &response, "/versions", nil, headers)
	return
}
