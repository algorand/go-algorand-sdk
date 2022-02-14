package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountAssetsParams defines parameters for LookupAccountAssets.
type LookupAccountAssetsParams struct {
	// Asset ID
	AssetID uint64 `url:"asset-id,omitempty"`

	// Include all items including closed accounts, deleted applications, destroyed assets, opted-out asset holdings, and closed-out application localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Maximum number of results to return. There could be additional pages even if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// The next page of results. Use the next token provided by the previous results.
	Next string `url:"next,omitempty"`
}

type LookupAccountAssets struct {
	c *Client

	accountID string

	p LookupAccountAssetsParams
}

// AssetID sets the asset ID filter.
func (s *LookupAccountAssets) AssetID(assetID uint64) *LookupAccountAssets {
	s.p.AssetID = assetID
	return s
}

// IncludeAll sets whether deleted assets will be requested.
func (s *LookupAccountAssets) IncludeAll(includeAll bool) *LookupAccountAssets {
	s.p.IncludeAll = includeAll
	return s
}

// Limit sets the limit for the number of returned assets.
func (s *LookupAccountAssets) Limit(limit uint64) *LookupAccountAssets {
	s.p.Limit = limit
	return s
}

// Next sets the next token for pagination. Use the next token provided by the previous
// results.
func (s *LookupAccountAssets) Next(next string) *LookupAccountAssets {
	s.p.Next = next
	return s
}

// Do performs the HTTP request.
func (s *LookupAccountAssets) Do(ctx context.Context, headers ...*common.Header) (response models.AssetHoldingsResponse, err error) {
	err = s.c.get(
		ctx, &response, fmt.Sprintf("/v2/accounts/%s/assets", s.accountID),
		s.p, headers)
	return
}
