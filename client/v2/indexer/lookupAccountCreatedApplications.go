package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// LookupAccountCreatedApplicationsParams contains all of the query parameters for url serialization.
type LookupAccountCreatedApplicationsParams struct {

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

// LookupAccountCreatedApplications lookup an account's created application
// parameters, optionally for a specific ID.
type LookupAccountCreatedApplications struct {
	c *Client

	accountId string

	p LookupAccountCreatedApplicationsParams
}

// ApplicationID application ID
func (s *LookupAccountCreatedApplications) ApplicationID(ApplicationID uint64) *LookupAccountCreatedApplications {
	s.p.ApplicationID = ApplicationID

	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupAccountCreatedApplications) IncludeAll(IncludeAll bool) *LookupAccountCreatedApplications {
	s.p.IncludeAll = IncludeAll

	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *LookupAccountCreatedApplications) Limit(Limit uint64) *LookupAccountCreatedApplications {
	s.p.Limit = Limit

	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *LookupAccountCreatedApplications) Next(Next string) *LookupAccountCreatedApplications {
	s.p.Next = Next

	return s
}

// Do performs the HTTP request
func (s *LookupAccountCreatedApplications) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationsResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s/created-applications", common.EscapeParams(s.accountId)...), s.p, headers)
	return
}
