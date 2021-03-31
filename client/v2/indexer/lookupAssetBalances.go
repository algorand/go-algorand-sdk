package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAssetBalancesParams contains all of the query parameters for url serialization.
type LookupAssetBalancesParams struct {

	// CurrencyGreaterThan results should have an amount greater than this value.
	// MicroAlgos are the default currency unless an asset-id is provided, in which
	// case the asset will be used.
	CurrencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// CurrencyLessThan results should have an amount less than this value. MicroAlgos
	// are the default currency unless an asset-id is provided, in which case the asset
	// will be used.
	CurrencyLessThan uint64 `url:"currency-less-than,omitempty"`

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// NextToken the next page of results. Use the next token provided by the previous
	// results.
	NextToken string `url:"next,omitempty"`

	// Round include results for the specified round.
	Round uint64 `url:"round,omitempty"`
}

// LookupAssetBalances lookup the list of accounts who hold this asset
type LookupAssetBalances struct {
	c *Client

	assetId uint64

	p LookupAssetBalancesParams
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *LookupAssetBalances) CurrencyGreaterThan(CurrencyGreaterThan uint64) *LookupAssetBalances {
	s.p.CurrencyGreaterThan = CurrencyGreaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
func (s *LookupAssetBalances) CurrencyLessThan(CurrencyLessThan uint64) *LookupAssetBalances {
	s.p.CurrencyLessThan = CurrencyLessThan
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupAssetBalances) IncludeAll(IncludeAll bool) *LookupAssetBalances {
	s.p.IncludeAll = IncludeAll
	return s
}

// Limit maximum number of results to return.
func (s *LookupAssetBalances) Limit(Limit uint64) *LookupAssetBalances {
	s.p.Limit = Limit
	return s
}

// NextToken the next page of results. Use the next token provided by the previous
// results.
func (s *LookupAssetBalances) NextToken(NextToken string) *LookupAssetBalances {
	s.p.NextToken = NextToken
	return s
}

// Round include results for the specified round.
func (s *LookupAssetBalances) Round(Round uint64) *LookupAssetBalances {
	s.p.Round = Round
	return s
}

// Do performs the HTTP request
func (s *LookupAssetBalances) Do(ctx context.Context, headers ...*common.Header) (response models.AssetBalancesResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/assets/%v/balances", s.assetId), s.p, headers)
	return
}
