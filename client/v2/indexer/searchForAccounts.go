package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type searchForAccountsParams struct {

	// applicationId application ID
	applicationId uint64 `url:"application-id,omitempty"`

	// assetID asset ID
	assetID uint64 `url:"asset-id,omitempty"`

	// authAddress include accounts configured to use this spending key.
	authAddress string `url:"auth-addr,omitempty"`

	// currencyGreaterThan results should have an amount greater than this value.
	// MicroAlgos are the default currency unless an asset-id is provided, in which
	// case the asset will be used.
	currencyGreaterThan uint64 `url:"currency-greater-than,omitempty"`

	// currencyLessThan results should have an amount less than this value. MicroAlgos
	// are the default currency unless an asset-id is provided, in which case the asset
	// will be used.
	currencyLessThan uint64 `url:"currency-less-than,omitempty"`

	// includeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	includeAll bool `url:"include-all,omitempty"`

	// limit maximum number of results to return.
	limit uint64 `url:"limit,omitempty"`

	// nextToken the next page of results. Use the next token provided by the previous
	// results.
	nextToken string `url:"next,omitempty"`

	// round include results for the specified round. For performance reasons, this
	// parameter may be disabled on some configurations.
	round uint64 `url:"round,omitempty"`
}

type SearchAccounts struct {
	c *Client

	p searchForAccountsParams
}

// ApplicationId application ID
func (s *SearchAccounts) ApplicationId(applicationId uint64) *SearchAccounts {
	s.p.applicationId = applicationId
	return s
}

// AssetID asset ID
func (s *SearchAccounts) AssetID(assetID uint64) *SearchAccounts {
	s.p.assetID = assetID
	return s
}

// AuthAddress include accounts configured to use this spending key.
func (s *SearchAccounts) AuthAddress(authAddress string) *SearchAccounts {
	s.p.authAddress = authAddress
	return s
}

// CurrencyGreaterThan results should have an amount greater than this value.
// MicroAlgos are the default currency unless an asset-id is provided, in which
// case the asset will be used.
func (s *SearchAccounts) CurrencyGreaterThan(currencyGreaterThan uint64) *SearchAccounts {
	s.p.currencyGreaterThan = currencyGreaterThan
	return s
}

// CurrencyLessThan results should have an amount less than this value. MicroAlgos
// are the default currency unless an asset-id is provided, in which case the asset
// will be used.
func (s *SearchAccounts) CurrencyLessThan(currencyLessThan uint64) *SearchAccounts {
	s.p.currencyLessThan = currencyLessThan
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *SearchAccounts) IncludeAll(includeAll bool) *SearchAccounts {
	s.p.includeAll = includeAll
	return s
}

// Limit maximum number of results to return.
func (s *SearchAccounts) Limit(limit uint64) *SearchAccounts {
	s.p.limit = limit
	return s
}

// NextToken the next page of results. Use the next token provided by the previous
// results.
func (s *SearchAccounts) NextToken(nextToken string) *SearchAccounts {
	s.p.nextToken = nextToken
	return s
}

// Round include results for the specified round. For performance reasons, this
// parameter may be disabled on some configurations.
func (s *SearchAccounts) Round(round uint64) *SearchAccounts {
	s.p.round = round
	return s
}

func (s *SearchAccounts) Do(ctx context.Context, headers ...*common.Header) (response models.AccountsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/accounts", s.p, headers)
	return
}
