package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAssetBalances is used to lookup asset balances
type LookupAssetBalances struct {
	c     *Client
	index uint64
	p     models.LookupAssetBalancesParams
}

func (s *LookupAssetBalances) NextToken(nextToken string) *LookupAssetBalances {
	s.p.NextToken = nextToken
	return s
}

func (s *LookupAssetBalances) Limit(lim uint64) *LookupAssetBalances {
	s.p.Limit = lim
	return s
}

func (s *LookupAssetBalances) Round(rnd uint64) *LookupAssetBalances {
	s.p.Round = rnd
	return s
}

func (s *LookupAssetBalances) CurrencyGreaterThan(greaterThan uint64) *LookupAssetBalances {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *LookupAssetBalances) CurrencyLessThan(lessThan uint64) *LookupAssetBalances {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *LookupAssetBalances) Do(ctx context.Context, headers ...*common.Header) (response models.AssetBalancesResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/assets/%d/balances", s.index), s.p, headers)
	return
}
