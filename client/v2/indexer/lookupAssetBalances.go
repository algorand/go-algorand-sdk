package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /v2/assets/{asset-id}/balances
 * Lookup the list of accounts who hold this asset
 */
type LookupAssetBalances struct {
	c       *Client
	p       models.LookupAssetBalancesParams
	assetId uint64
}

/**
 * Results should have an amount greater than this value. MicroAlgos are the
 * default currency unless an asset-id is provided, in which case the asset will be
 * used.
 */
func (s *LookupAssetBalances) CurrencyGreaterThan(currencyGreaterThan uint64) *LookupAssetBalances {
	s.p.CurrencyGreaterThan = currencyGreaterThan
	return s
}

/**
 * Results should have an amount less than this value. MicroAlgos are the default
 * currency unless an asset-id is provided, in which case the asset will be used.
 */
func (s *LookupAssetBalances) CurrencyLessThan(currencyLessThan uint64) *LookupAssetBalances {
	s.p.CurrencyLessThan = currencyLessThan
	return s
}

/**
 * Maximum number of results to return.
 */
func (s *LookupAssetBalances) Limit(limit uint64) *LookupAssetBalances {
	s.p.Limit = limit
	return s
}

/**
 * The next page of results. Use the next token provided by the previous results.
 */
func (s *LookupAssetBalances) NextToken(next string) *LookupAssetBalances {
	s.p.Next = next
	return s
}

/**
 * Include results for the specified round.
 */
func (s *LookupAssetBalances) Round(round uint64) *LookupAssetBalances {
	s.p.Round = round
	return s
}

func (s *LookupAssetBalances) Do(ctx context.Context,
	headers ...*common.Header) (response models.AssetBalancesResponse, err error) {
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/assets/%d/balances", s.assetId), s.p, headers)
	return
}
