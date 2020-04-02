package indexer

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAssetBalancesService is used to lookup asset balances
type LookupAssetBalancesService struct {
	c     *Client
	index uint64
	p     models.LookupAssetBalancesParams
}

func (s *LookupAssetBalancesService) Limit(lim uint64) *LookupAssetBalancesService {
	s.p.Limit = lim
	return s
}

func (s *LookupAssetBalancesService) AfterAddress(after string) *LookupAssetBalancesService {
	s.p.AfterAddress = after
	return s
}

func (s *LookupAssetBalancesService) Round(rnd uint64) *LookupAssetBalancesService {
	s.p.Round = rnd
	return s
}

func (s *LookupAssetBalancesService) CurrencyGreaterThan(greaterThan uint64) *LookupAssetBalancesService {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *LookupAssetBalancesService) CurrencyLessThan(lessThan uint64) *LookupAssetBalancesService {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *LookupAssetBalancesService) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, holders []models.MiniAssetHolding, err error) {
	var response models.AssetBalancesResponse
	err = s.c.get(ctx, &response, fmt.Sprintf("/assets/%d/balances", s.index), s.p, headers)
	validRound = response.CurrentRound
	holders = response.Balances
	return
}
