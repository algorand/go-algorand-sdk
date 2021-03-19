package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type GetApplicationByID struct {
	c *Client

	applicationId uint64
}

func (s *GetApplicationByID) Do(ctx context.Context, headers ...*common.Header) (response models.Application, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/applications/%v", s.applicationId), nil, headers)
	return
}
