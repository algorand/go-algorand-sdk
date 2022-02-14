package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// AccountApplicationInformationParams defines parameters for AccountApplicationInformation.
type AccountApplicationInformationParams struct {
	// Configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

type AccountApplicationInformation struct {
	c *Client

	address string
	applicationID uint64

	p AccountApplicationInformationParams
}

// Do performs the HTTP request.
func (s *AccountApplicationInformation) Do(ctx context.Context, headers ...*common.Header) (response models.AccountApplicationResponse, err error) {
	err = s.c.get(
		ctx, &response,
		fmt.Sprintf("/v2/accounts/%s/application/%d", s.address, s.applicationID),
		s.p, headers)
	return
}
