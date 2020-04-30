package indexer

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type SearchAccounts struct {
	c *Client
	p models.SearchAccountsParams
}

func (s *SearchAccounts) Next(nextToken string) *SearchAccounts {
	s.p.NextToken = nextToken
	return s
}

func (s *SearchAccounts) AssetID(assetID uint64) *SearchAccounts {
	s.p.AssetId = assetID
	return s
}

func (s *SearchAccounts) Limit(limit uint64) *SearchAccounts {
	s.p.Limit = limit
	return s
}
func (s *SearchAccounts) CurrencyGreaterThan(greaterThan uint64) *SearchAccounts {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

func (s *SearchAccounts) CurrencyLessThan(lessThan uint64) *SearchAccounts {
	s.p.CurrencyLessThan = lessThan
	return s
}

func (s *SearchAccounts) AfterAddress(after string) *SearchAccounts {
	s.p.AfterAddress = after
	return s
}

func (s *SearchAccounts) Round(round uint64) *SearchAccounts {
	s.p.Round = round
	return s
}

func (s *SearchAccounts) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result []models.Account, err error) {
	var response models.AccountsResponse
	err = s.c.get(ctx, &response, "/accounts", s.p, headers)
	validRound = response.CurrentRound
	result = response.Accounts
	return
}
