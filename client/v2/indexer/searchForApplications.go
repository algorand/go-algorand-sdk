package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// SearchForApplicationsParams contains all of the query parameters for url serialization.
type SearchForApplicationsParams struct {

	// ApplicationId application ID
	ApplicationId uint64 `url:"application-id,omitempty"`

	// Creator filter just applications with the given creator address.
	Creator string `url:"creator,omitempty"`

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

// SearchForApplications search for applications
type SearchForApplications struct {
	c *Client

	p SearchForApplicationsParams
}

// ApplicationId application ID
func (s *SearchForApplications) ApplicationId(ApplicationId uint64) *SearchForApplications {
	s.p.ApplicationId = ApplicationId
	return s
}

// Creator filter just applications with the given creator address.
func (s *SearchForApplications) Creator(Creator string) *SearchForApplications {
	s.p.Creator = Creator
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *SearchForApplications) IncludeAll(IncludeAll bool) *SearchForApplications {
	s.p.IncludeAll = IncludeAll
	return s
}

// Limit maximum number of results to return. There could be additional pages even
// if the limit is not reached.
func (s *SearchForApplications) Limit(Limit uint64) *SearchForApplications {
	s.p.Limit = Limit
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *SearchForApplications) Next(Next string) *SearchForApplications {
	s.p.Next = Next
	return s
}

// Do performs the HTTP request
func (s *SearchForApplications) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/applications", s.p, headers)
	return
}
