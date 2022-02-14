package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupAccountAppLocalStatesParams defines parameters for LookupAccountAppLocalStates.
type LookupAccountAppLocalStatesParams struct {
	// Application ID
	ApplicationID uint64 `url:"application-id,omitempty"`

	// Include all items including closed accounts, deleted applications, destroyed assets, opted-out asset holdings, and closed-out application localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Maximum number of results to return. There could be additional pages even if the limit is not reached.
	Limit uint64 `url:"limit,omitempty"`

	// The next page of results. Use the next token provided by the previous results.
	Next string `url:"next,omitempty"`
}

type LookupAccountAppLocalStates struct {
	c *Client

	accountID string

	p LookupAccountAppLocalStatesParams
}

// AssetID sets the application ID filter.
func (s *LookupAccountAppLocalStates) ApplicationID(applicationID uint64) *LookupAccountAppLocalStates {
	s.p.ApplicationID = applicationID
	return s
}

// IncludeAll sets whether deleted applications will be requested.
func (s *LookupAccountAppLocalStates) IncludeAll(includeAll bool) *LookupAccountAppLocalStates {
	s.p.IncludeAll = includeAll
	return s
}

// Limit sets the limit for the number of returned assets.
func (s *LookupAccountAppLocalStates) Limit(limit uint64) *LookupAccountAppLocalStates {
	s.p.Limit = limit
	return s
}

// Next sets the next token for pagination. Use the next token provided by the previous
// results.
func (s *LookupAccountAppLocalStates) Next(next string) *LookupAccountAppLocalStates {
	s.p.Next = next
	return s
}

// Do performs the HTTP request.
func (s *LookupAccountAppLocalStates) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationLocalStatesResponse, err error) {
	err = s.c.get(
		ctx, &response, fmt.Sprintf("/v2/accounts/%s/apps-local-state", s.accountID),
		s.p, headers)
	return
}
