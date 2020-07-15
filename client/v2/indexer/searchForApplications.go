package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// SearchForApplications /v2/applications
// Search for applications
type SearchForApplications struct {
	c *Client
	p models.SearchForApplicationsParams
}

// ApplicationId application ID
func (s *SearchForApplications) ApplicationId(applicationId uint64) *SearchForApplications {
	s.p.ApplicationId = applicationId
	return s
}

// Do performs HTTP request
func (s *SearchForApplications) Do(ctx context.Context,
	headers ...*common.Header) (response models.ApplicationsResponse, err error) {
	err = s.c.get(ctx, &response,
		"/v2/applications", s.p, headers)
	return
}
