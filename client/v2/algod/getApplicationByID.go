package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetApplicationByID given a application ID, it returns application information
// including creator, approval and clear programs, global and local schemas, and
// global state.
type GetApplicationByID struct {
	c *Client

	applicationId uint64
}

// Do performs the HTTP request
func (s *GetApplicationByID) Do(ctx context.Context, headers ...*common.Header) (response models.Application, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s", common.EscapeParams(s.applicationId)...), nil, headers)
	return
}
