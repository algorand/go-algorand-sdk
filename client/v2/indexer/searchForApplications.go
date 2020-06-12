package indexer

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /v2/applications
 * Search for applications
 */
type SearchForApplications struct {
	c *Client
	p models.SearchForApplicationsParams
}

/**
 * Application ID
 */
func (s *SearchForApplications) ApplicationId(applicationId uint64) *SearchForApplications {
	s.p.ApplicationId = applicationId
	return s
}

/**
 * Include results for the specified round.
 */
func (s *SearchForApplications) Round(round uint64) *SearchForApplications {
	s.p.Round = round
	return s
}

func (s *SearchForApplications) Do(ctx context.Context,
	headers ...*common.Header) (response models.ApplicationsResponse, err error) {
	err = s.c.get(ctx, &response,
		"/v2/applications", s.p, headers)
	return
}
