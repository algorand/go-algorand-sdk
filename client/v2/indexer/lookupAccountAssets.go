package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountAssetsParams contains all of the query parameters for url serialization.
type LookupAccountAssetsParams struct {

	// AssetID asset ID
	AssetID uint64 `url:"asset-id,omitempty"`

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Limit maximum number of results to return. There could be additional pages even
	// if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`
}

// LookupAccountAssets lookup an account's asset holdings, optionally for a
// specific ID.
type LookupAccountAssets struct {
	c *Client

	accountId string

	p LookupAccountAssetsParams
}

// AssetID asset ID
func (s *LookupAccountAssets) AssetID(AssetID uint64) *LookupAccountAssets {
	s.p.AssetID = AssetID
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupAccountAssets) IncludeAll(IncludeAll bool) *LookupAccountAssets {
	s.p.IncludeAll = IncludeAll
	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *LookupAccountAssets) Limit(Limit uint64) *LookupAccountAssets {
	s.p.Limit = Limit
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *LookupAccountAssets) Next(Next string) *LookupAccountAssets {
	s.p.Next = Next
	return s
}

// Do performs the HTTP request
func (s *LookupAccountAssets) Do(ctx context.Context, headers ...*common.Header) (response models.AssetHoldingsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s/assets", common.EscapeParams(s.accountId)...), s.p, headers)
	return
}
