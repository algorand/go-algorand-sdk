package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type searchForApplicationsParams struct {

	// applicationId application ID
	applicationId uint64 `url:"application-id,omitempty"`

	// includeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	includeAll bool `url:"include-all,omitempty"`

	// limit maximum number of results to return.
	limit uint64 `url:"limit,omitempty"`

	// next the next page of results. Use the next token provided by the previous
	// results.
	next string `url:"next,omitempty"`
}

type SearchForApplications struct {
	c *Client

	p searchForApplicationsParams
}

// ApplicationId application ID
func (s *SearchForApplications) ApplicationId(applicationId uint64) *SearchForApplications {
	s.p.applicationId = applicationId
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *SearchForApplications) IncludeAll(includeAll bool) *SearchForApplications {
	s.p.includeAll = includeAll
	return s
}

// Limit maximum number of results to return.
func (s *SearchForApplications) Limit(limit uint64) *SearchForApplications {
	s.p.limit = limit
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *SearchForApplications) Next(next string) *SearchForApplications {
	s.p.next = next
	return s
}

func (s *SearchForApplications) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/applications", s.p, headers)
	return
}
