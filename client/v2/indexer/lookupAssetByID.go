package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupAssetByID struct {
	c     *Client
	index uint64
}

func (s *LookupAssetByID) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result models.Asset, err error) {
	response := models.LookupAssetByIDResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/assets/%d", s.index), nil, headers)
	validRound = response.CurrentRound
	result = response.Asset
	return
}
