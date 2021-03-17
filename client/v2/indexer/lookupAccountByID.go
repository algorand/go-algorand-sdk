package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type LookupAccountByIDParams struct {

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Round include results for the specified round.
	Round uint64 `url:"round,omitempty"`
}

type LookupAccountByID struct {
	c *Client

	accountId string

	p LookupAccountByIDParams
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

func (s *LookupAccountByID) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result models.Account, err error) {
	response := models.LookupAccountByIDResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s", s.account), s.p, headers)
	validRound = response.CurrentRound
	result = response.Account
	return
}
