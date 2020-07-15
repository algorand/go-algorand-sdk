package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// LookupApplicationByID /v2/applications/{application-id}
// Lookup application.
type LookupApplicationByID struct {
	c             *Client
	applicationId uint64
}

// Do performs HTTP request
func (s *LookupApplicationByID) Do(ctx context.Context,
	headers ...*common.Header) (response models.ApplicationResponse, err error) {
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/applications/%d", s.applicationId), nil, headers)
	return
}
