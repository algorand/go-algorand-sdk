package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// SearchAccountsParams contains all of the query parameters for url serialization.
type SearchAccountsParams struct {

	// ApplicationId application ID
	ApplicationId uint64 `url:"application-id,omitempty"`

	// AssetID asset ID
	AssetID uint64 `url:"asset-id,omitempty"`

	// AuthAddress include accounts configured to use this spending key.
	AuthAddress string `url:"auth-addr,omitempty"`

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

	// Round include results for the specified round. For performance reasons, this
	// parameter may be disabled on some configurations.
	Round uint64 `url:"round,omitempty"`
}

// SearchAccounts search for accounts.
type SearchAccounts struct {
	c *Client

	p SearchAccountsParams
}

// ApplicationId application ID
func (s *SearchAccounts) ApplicationId(ApplicationId uint64) *SearchAccounts {
	s.p.ApplicationId = ApplicationId
	return s
}

// AssetID asset ID
func (s *SearchAccounts) AssetID(AssetID uint64) *SearchAccounts {
	s.p.AssetID = AssetID
	return s
}

// AuthAddress include accounts configured to use this spending key.
func (s *SearchAccounts) AuthAddress(AuthAddress string) *SearchAccounts {
	s.p.AuthAddress = AuthAddress
	return s
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *SearchAccounts) CurrencyGreaterThan(CurrencyGreaterThan uint64) *SearchAccounts {
	s.p.CurrencyGreaterThan = CurrencyGreaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
func (s *SearchAccounts) CurrencyLessThan(CurrencyLessThan uint64) *SearchAccounts {
	s.p.CurrencyLessThan = CurrencyLessThan
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *SearchAccounts) IncludeAll(IncludeAll bool) *SearchAccounts {
	s.p.IncludeAll = IncludeAll
	return s
}

// Limit maximum number of results to return.
func (s *SearchAccounts) Limit(Limit uint64) *SearchAccounts {
	s.p.Limit = Limit
	return s
}

// NextToken the next page of results. Use the next token provided by the previous
// results.
func (s *SearchAccounts) NextToken(NextToken string) *SearchAccounts {
	s.p.NextToken = NextToken
	return s
}

// Round include results for the specified round. For performance reasons, this
// parameter may be disabled on some configurations.
func (s *SearchAccounts) Round(Round uint64) *SearchAccounts {
	s.p.Round = Round
	return s
}

// Do performs the HTTP request
func (s *SearchAccounts) Do(ctx context.Context, headers ...*common.Header) (response models.AccountsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/accounts", s.p, headers)
	return
}
