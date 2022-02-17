package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountByIDParams contains all of the query parameters for url serialization.
type LookupAccountByIDParams struct {

	// Exclude exclude additional items such as asset holdings, application local data
	// stored for this account, asset parameters created by this account, and
	// application parameters created by this account.
	Exclude []string `url:"exclude,omitempty,comma"`

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Round include results for the specified round.
	Round uint64 `url:"round,omitempty"`
}

// LookupAccountByID lookup account information.
type LookupAccountByID struct {
	c *Client

	accountId string

	p LookupAccountByIDParams
}

// Exclude exclude additional items such as asset holdings, application local data
// stored for this account, asset parameters created by this account, and
// application parameters created by this account.
func (s *LookupAccountByID) Exclude(Exclude []string) *LookupAccountByID {
	s.p.Exclude = Exclude
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupAccountByID) IncludeAll(IncludeAll bool) *LookupAccountByID {
	s.p.IncludeAll = IncludeAll
	return s
}

// Round include results for the specified round.
func (s *LookupAccountByID) Round(Round uint64) *LookupAccountByID {
	s.p.Round = Round
	return s
}

// Do performs the HTTP request
func (s *LookupAccountByID) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result models.Account, err error) {
	response := models.AccountResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s", s.accountId), s.p, headers)
	validRound = response.CurrentRound
	result = response.Account
	return
}
