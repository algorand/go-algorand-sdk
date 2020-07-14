package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetAssetByID /v2/assets/{asset-id}
// Given a asset id, it returns asset information including creator, name, total
// supply and special addresses.
type GetAssetByID struct {
	c       *Client
	assetId uint64
}

// Do performs HTTP request
func (s *GetAssetByID) Do(ctx context.Context,
	headers ...*common.Header) (response models.Asset, err error) {
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/assets/%d", s.assetId), nil, headers)
	return
}
