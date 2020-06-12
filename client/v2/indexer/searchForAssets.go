package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /v2/assets
 * Search for assets.
 */
type SearchForAssets struct {
	c *Client
	p models.SearchForAssetsParams
}

/**
 * Asset ID
 */
func (s *SearchForAssets) AssetID(assetId uint64) *SearchForAssets {
	s.p.AssetId = assetId
	return s
}

/**
 * Filter just assets with the given creator address.
 */
func (s *SearchForAssets) Creator(creator string) *SearchForAssets {
	s.p.Creator = creator
	return s
}

/**
 * Maximum number of results to return.
 */
func (s *SearchForAssets) Limit(limit uint64) *SearchForAssets {
	s.p.Limit = limit
	return s
}

/**
 * Filter just assets with the given name.
 */
func (s *SearchForAssets) Name(name string) *SearchForAssets {
	s.p.Name = name
	return s
}

/**
 * The next page of results. Use the next token provided by the previous results.
 */
func (s *SearchForAssets) NextToken(next string) *SearchForAssets {
	s.p.Next = next
	return s
}

/**
 * Filter just assets with the given unit.
 */
func (s *SearchForAssets) Unit(unit string) *SearchForAssets {
	s.p.Unit = unit
	return s
}

func (s *SearchForAssets) Do(ctx context.Context,
	headers ...*common.Header) (validRound uint64, result []models.Asset, err error) {
	var response models.AssetsResponse
	err = s.c.get(ctx, &response,
		"/v2/assets", s.p, headers)
	validRound = response.CurrentRound
	result = response.Assets
	return
}
