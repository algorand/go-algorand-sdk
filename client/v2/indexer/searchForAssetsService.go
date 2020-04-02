package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type SearchForAssetsService struct {
	c *Client
	p models.SearchForAssetsParams
}

func (s *SearchForAssetsService) Limit(lim uint64) *SearchForAssetsService {
	s.p.Limit = lim
	return s
}

func (s *SearchForAssetsService) Creator(creator string) *SearchForAssetsService {
	s.p.Creator = creator
	return s
}

func (s *SearchForAssetsService) Name(name string) *SearchForAssetsService {
	s.p.Name = name
	return s
}

func (s *SearchForAssetsService) Unit(unit string) *SearchForAssetsService {
	s.p.Unit = unit
	return s
}

func (s *SearchForAssetsService) AssetID(id uint64) *SearchForAssetsService {
	s.p.AssetId = id
	return s
}

func (s *SearchForAssetsService) AfterAsset(after uint64) *SearchForAssetsService {
	s.p.AfterAsset = after
	return s
}

func (s *SearchForAssetsService) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result []models.Asset, err error) {
	var response models.AssetsResponse
	err = s.c.get(ctx, &response, "/assets", s.p, headers)
	validRound = response.CurrentRound
	result = response.Assets
	return
}
