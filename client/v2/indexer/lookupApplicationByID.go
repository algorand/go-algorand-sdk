package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupApplicationByIDParams contains all of the query parameters for url serialization.
type LookupApplicationByIDParams struct {

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`
}

// LookupApplicationByID lookup application.
type LookupApplicationByID struct {
	c *Client

	applicationId uint64

	p LookupApplicationByIDParams
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupApplicationByID) IncludeAll(IncludeAll bool) *LookupApplicationByID {
	s.p.IncludeAll = IncludeAll
	return s
}

// Do performs the HTTP request
func (s *LookupApplicationByID) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
