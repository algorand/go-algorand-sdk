package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type lookupApplicationByIDParams struct {

	// includeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	includeAll bool `url:"include-all,omitempty"`
}

type LookupApplicationByID struct {
	c *Client

	applicationId uint64

	p lookupApplicationByIDParams
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupApplicationByID) IncludeAll(includeAll bool) *LookupApplicationByID {
	s.p.includeAll = includeAll
	return s
}

func (s *LookupApplicationByID) Do(ctx context.Context, headers ...*common.Header) (response models.ApplicationResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%v", s.applicationId), s.p, headers)
	return
}
