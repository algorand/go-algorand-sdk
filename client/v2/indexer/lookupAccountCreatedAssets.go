package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountCreatedAssetsParams defines parameters for LookupAccountCreatedAssets.
type LookupAccountCreatedAssetsParams struct {
	// Asset ID
	AssetID uint64 `url:"asset-id,omitempty"`

	// Include all items including closed accounts, deleted applications, destroyed assets, opted-out asset holdings, and closed-out application localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Maximum number of results to return. There could be additional pages even if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// The next page of results. Use the next token provided by the previous results.
	Next string `url:"next,omitempty"`
}

type LookupAccountCreatedAssets struct {
	c *Client

	accountID string

	p LookupAccountCreatedAssetsParams
}

// AssetID sets the asset ID filter.
func (s *LookupAccountCreatedAssets) AssetID(assetID uint64) *LookupAccountCreatedAssets {
	s.p.AssetID = assetID
	return s
}

// IncludeAll sets whether deleted assets will be requested.
func (s *LookupAccountCreatedAssets) IncludeAll(includeAll bool) *LookupAccountCreatedAssets {
	s.p.IncludeAll = includeAll
	return s
}

// Limit sets the limit for the number of returned assets.
func (s *LookupAccountCreatedAssets) Limit(limit uint64) *LookupAccountCreatedAssets {
	s.p.Limit = limit
	return s
}

// Next sets the next token for pagination. Use the next token provided by the previous
// results.
func (s *LookupAccountCreatedAssets) Next(next string) *LookupAccountCreatedAssets {
	s.p.Next = next
	return s
}

// Do performs the HTTP request.
func (s *LookupAccountCreatedAssets) Do(ctx context.Context, headers ...*common.Header) (response models.AssetsResponse, err error) {
	err = s.c.get(
		ctx, &response, fmt.Sprintf("/v2/accounts/%s/created-assets", s.accountID),
		s.p, headers)
	return
}
