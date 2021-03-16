package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type searchForApplicationsParams struct {

	// ApplicationId application ID
	ApplicationId uint64 `url:"application-id,omitempty"`

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`
}

type SearchForApplications struct {
	c *Client

	p searchForApplicationsParams
}

// ApplicationId application ID
func (s *SearchForApplications) ApplicationId(applicationId uint64) *SearchForApplications {
	s.p.ApplicationId = applicationId
	return s
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *SearchForApplications) IncludeAll(includeAll bool) *SearchForApplications {
	s.p.IncludeAll = includeAll
	return s
}

// Limit maximum number of results to return.
func (s *SearchForApplications) Limit(limit uint64) *SearchForApplications {
	s.p.Limit = limit
	return s
}

// Next the next page of results. Use the next token provided by the previous
// results.
func (s *SearchForApplications) Next(next string) *SearchForApplications {
	s.p.Next = next
	return s
}

func (s *SearchForApplications) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationsResponse, err error) {
	err = s.c.get(ctx, &response, "/v2/applications", s.p, headers)
	return
}
