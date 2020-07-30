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

// ApplicationId application ID
func (s *SearchAccounts) ApplicationId(applicationId uint64) *SearchAccounts {
	s.p.ApplicationId = applicationId
	return s
}

// NextToken the next page of results. Use the next token provided by the previous
// results.
func (s *SearchAccounts) NextToken(nextToken string) *SearchAccounts {
	s.p.NextToken = nextToken
	return s
}

func (s *SearchAccounts) AssetID(assetID uint64) *SearchAccounts {
	s.p.AssetId = assetID
	return s
}

// Limit maximum number of results to return.
func (s *SearchAccounts) Limit(limit uint64) *SearchAccounts {
	s.p.Limit = limit
	return s
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *SearchAccounts) CurrencyGreaterThan(greaterThan uint64) *SearchAccounts {
	s.p.CurrencyGreaterThan = greaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
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

func (s *SearchAccounts) AuthAddress(authAddr string) *SearchAccounts {
	s.p.AuthAddr = authAddr
	return s
}

func (s *SearchAccounts) Do(ctx context.Context, headers ...*common.Header) (response models.AccountsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/accounts", s.p, headers)
	return
}
