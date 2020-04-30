package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type SearchForAssets struct {
	c *Client
	p models.SearchForAssetsParams
}

func (s *SearchForAssets) NextToken(nextToken string) *SearchForAssets {
	s.p.NextToken = nextToken
	return s
}

func (s *SearchForAssets) Limit(lim uint64) *SearchForAssets {
	s.p.Limit = lim
	return s
}

func (s *SearchForAssets) Creator(creator string) *SearchForAssets {
	s.p.Creator = creator
	return s
}

func (s *SearchForAssets) Name(name string) *SearchForAssets {
	s.p.Name = name
	return s
}

func (s *SearchForAssets) Unit(unit string) *SearchForAssets {
	s.p.Unit = unit
	return s
}

func (s *SearchForAssets) AssetID(id uint64) *SearchForAssets {
	s.p.AssetId = id
	return s
}

func (s *SearchForAssets) AfterAsset(after uint64) *SearchForAssets {
	s.p.AfterAsset = after
	return s
}

func (s *SearchForAssets) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result []models.Asset, err error) {
	var response models.AssetsResponse
	err = s.c.get(ctx, &response, "/assets", s.p, headers)
	validRound = response.CurrentRound
	result = response.Assets
	return
}
