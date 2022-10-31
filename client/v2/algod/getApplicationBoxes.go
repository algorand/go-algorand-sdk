package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetApplicationBoxesParams contains all of the query parameters for url serialization.
type GetApplicationBoxesParams struct {

	// Max max number of box names to return. If max is not set, or max == 0, returns
	// all box-names.
	Max uint64 `url:"max,omitempty"`
}

// GetApplicationBoxes given an application ID, return all Box names. No particular
// ordering is guaranteed. Request fails when client or server-side configured
// limits prevent returning all Box names.
type GetApplicationBoxes struct {
	c *Client

	applicationId uint64

	p GetApplicationBoxesParams
}

// Max max number of box names to return. If max is not set, or max == 0, returns
// all box-names.
func (s *GetApplicationBoxes) Max(Max uint64) *GetApplicationBoxes {
	s.p.Max = Max

	return s
}

// Do performs the HTTP request
func (s *GetApplicationBoxes) Do(ctx context.Context, headers ...*common.Header) (response models.BoxesResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/boxes", common.EscapeParams(s.applicationId)...), s.p, headers)
	return
}
