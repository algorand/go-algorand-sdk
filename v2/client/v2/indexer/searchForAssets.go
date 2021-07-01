package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// SearchForAssetsParams contains all of the query parameters for url serialization.
type SearchForAssetsParams struct {

	// AssetID asset ID
	AssetID uint64 `url:"asset-id,omitempty"`

	// Creator filter just assets with the given creator address.
	Creator string `url:"creator,omitempty"`

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Name filter just assets with the given name.
	Name string `url:"name,omitempty"`

	// NextToken the next page of results. Use the next token provided by the previous
	// results.
	NextToken string `url:"next,omitempty"`

	// Unit filter just assets with the given unit.
	Unit string `url:"unit,omitempty"`
}

// SearchForAssets search for assets.
type SearchForAssets struct {
	c *Client

	p SearchForAssetsParams
}

// AssetID asset ID
func (s *SearchForAssets) AssetID(AssetID uint64) *SearchForAssets {
	s.p.AssetID = AssetID
	return s
}

// Creator filter just assets with the given creator address.
func (s *SearchForAssets) Creator(Creator string) *SearchForAssets {
	s.p.Creator = Creator
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *SearchForAssets) IncludeAll(IncludeAll bool) *SearchForAssets {
	s.p.IncludeAll = IncludeAll
	return s
}

// Limit maximum number of results to return.
func (s *SearchForAssets) Limit(Limit uint64) *SearchForAssets {
	s.p.Limit = Limit
	return s
}

// Name filter just assets with the given name.
func (s *SearchForAssets) Name(Name string) *SearchForAssets {
	s.p.Name = Name
	return s
}

// NextToken the next page of results. Use the next token provided by the previous
// results.
func (s *SearchForAssets) NextToken(NextToken string) *SearchForAssets {
	s.p.NextToken = NextToken
	return s
}

// Unit filter just assets with the given unit.
func (s *SearchForAssets) Unit(Unit string) *SearchForAssets {
	s.p.Unit = Unit
	return s
}

// Do performs the HTTP request
func (s *SearchForAssets) Do(ctx context.Context, headers ...*common.Header) (response models.AssetsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/assets", s.p, headers)
	return
}
