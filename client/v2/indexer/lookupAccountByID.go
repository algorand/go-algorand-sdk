package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountByIDParams contains all of the query parameters for url serialization.
type LookupAccountByIDParams struct {
	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Round include results for the specified round.
	Round uint64 `url:"round,omitempty"`

	// Exclude additional items such as asset holdings, application local data stored for this account, asset parameters created by this account, and application parameters created by this account.
	Exclude []string `url:"exclude,omitempty,comma"`
}

// LookupAccountByID lookup account information.
type LookupAccountByID struct {
	c *Client

	accountId string

	p LookupAccountByIDParams
}

// IncludeAll sets whether deleted accounts and creatables will be requested.
func (s *LookupAccountByID) IncludeAll(IncludeAll bool) *LookupAccountByID {
	s.p.IncludeAll = IncludeAll
	return s
}

// Round include results for the specified round.
func (s *LookupAccountByID) Round(Round uint64) *LookupAccountByID {
	s.p.Round = Round
	return s
}

// Exclude sets which creatable types must be excluded from the result. `assets` is true
// if and only if asset holdings must be excluded. Similar logic applies to the other
// creatable types.
func (s *LookupAccountByID) Exclude(assets, createdAssets, appsLocalState, createdApps bool) *LookupAccountByID {
	if assets {
		s.p.Exclude = append(s.p.Exclude, "assets")
	}
	if createdAssets {
		s.p.Exclude = append(s.p.Exclude, "created-assets")
	}
	if appsLocalState {
		s.p.Exclude = append(s.p.Exclude, "apps-local-state")
	}
	if createdApps {
		s.p.Exclude = append(s.p.Exclude, "created-apps")
	}

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
