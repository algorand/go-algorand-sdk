package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

/**
 * /v2/applications/{application-id}
 * Lookup application.
 */
type LookupApplication struct {
	c             *Client
	p             models.LookupApplicationParams
	applicationId uint64
}

/**
 * Include results for the specified round.
 */
func (s *LookupApplication) Round(round uint64) *LookupApplication {
	s.p.Round = round
	return s
}

func (s *LookupApplication) Do(ctx context.Context,
	headers ...*common.Header) (response models.ApplicationResponse, err error) {
	err = s.c.get(ctx, &response,
		fmt.Sprintf("/v2/applications/%d", s.applicationId), s.p, headers)
	return
}
