package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type GetAssetByID struct {
	c *Client

	assetId uint64
}

func (s *GetAssetByID) Do(ctx context.Context, headers ...*common.Header) (response models.Asset, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/assets/%v", s.assetId), nil, headers)
	return
}
