package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountAppLocalStatesParams contains all of the query parameters for url serialization.
type LookupAccountAppLocalStatesParams struct {

	// ApplicationID application ID
	ApplicationID uint64 `url:"application-id,omitempty"`

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Limit maximum number of results to return. There could be additional pages even
	// if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`
}

// LookupAccountAppLocalStates lookup an account's asset holdings, optionally for a
// specific ID.
type LookupAccountAppLocalStates struct {
	c *Client

	accountId string

	p LookupAccountAppLocalStatesParams
}

// ApplicationID application ID
func (s *LookupAccountAppLocalStates) ApplicationID(ApplicationID uint64) *LookupAccountAppLocalStates {
	s.p.ApplicationID = ApplicationID
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupAccountAppLocalStates) IncludeAll(IncludeAll bool) *LookupAccountAppLocalStates {
	s.p.IncludeAll = IncludeAll
	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *LookupAccountAppLocalStates) Limit(Limit uint64) *LookupAccountAppLocalStates {
	s.p.Limit = Limit
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *LookupAccountAppLocalStates) Next(Next string) *LookupAccountAppLocalStates {
	s.p.Next = Next
	return s
}

// Do performs the HTTP request
func (s *LookupAccountAppLocalStates) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationLocalStatesResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%v/apps-local-state", s.accountId), s.p, headers)
	return
}
