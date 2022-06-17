package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// GetApplicationBoxByName given an application ID and box name, it returns the box
// name and value (each base64 encoded).
type GetApplicationBoxByName struct {
	c *Client

	applicationId uint64
	boxName       string
}

// Do performs the HTTP request
func (s *GetApplicationBoxByName) Do(ctx context.Context, headers ...*common.Header) (response models.Box, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%s/boxes/%s", common.EscapeParams(s.applicationId, s.boxName)...), nil, headers)
	return
}
