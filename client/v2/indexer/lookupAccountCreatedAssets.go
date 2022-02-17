package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountCreatedAssetsParams contains all of the query parameters for url serialization.
type LookupAccountCreatedAssetsParams struct {

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

// LookupAccountCreatedAssets lookup an account's created asset parameters,
// optionally for a specific ID.
type LookupAccountCreatedAssets struct {
	c *Client

	accountId string

	p LookupAccountCreatedAssetsParams
}

// AssetID asset ID
func (s *LookupAccountCreatedAssets) AssetID(AssetID uint64) *LookupAccountCreatedAssets {
	s.p.AssetID = AssetID
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupAccountCreatedAssets) IncludeAll(IncludeAll bool) *LookupAccountCreatedAssets {
	s.p.IncludeAll = IncludeAll
	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *LookupAccountCreatedAssets) Limit(Limit uint64) *LookupAccountCreatedAssets {
	s.p.Limit = Limit
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *LookupAccountCreatedAssets) Next(Next string) *LookupAccountCreatedAssets {
	s.p.Next = Next
	return s
}

// Do performs the HTTP request
func (s *LookupAccountCreatedAssets) Do(ctx context.Context, headers ...*common.Header) (response models.AssetsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%v/created-assets", s.accountId), s.p, headers)
	return
}
