package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type SearchAccountsService struct {
	c *Client
	p models.SearchAccountsParams
}

func (s *SearchAccountsService) AssetID(assetID uint64) *SearchAccountsService {
	s.p.AssetId = assetID
	return s
}

func (s *SearchAccountsService) Limit(limit uint64) *SearchAccountsService {
	s.p.Limit = limit
	return s
}
func (s *SearchAccountsService) CurrencyGreaterThan(greaterThan uint64) *SearchAccountsService {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *SearchAccountsService) CurrencyLessThan(lessThan uint64) *SearchAccountsService {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *SearchAccountsService) AfterAddress(after string) *SearchAccountsService {
	s.p.AfterAddress = after
	return s
}

func (s *SearchAccountsService) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result []models.Account, err error) {
	var response models.AccountsResponse
	err = s.c.get(ctx, &response, "/accounts", s.p, headers)
	validRound = response.CurrentRound
	result = response.Accounts
	return
}
