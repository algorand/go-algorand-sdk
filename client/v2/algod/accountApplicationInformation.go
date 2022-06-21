package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// AccountApplicationInformationParams contains all of the query parameters for url serialization.
type AccountApplicationInformationParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

// AccountApplicationInformation given a specific account public key and
// application ID, this call returns the account's application local state and
// global state (AppLocalState and AppParams, if either exists). Global state will
// only be returned if the provided address is the application's creator.
type AccountApplicationInformation struct {
	c *Client

	address       string
	applicationId uint64

	p AccountApplicationInformationParams
}

// Do performs the HTTP request
func (s *AccountApplicationInformation) Do(ctx context.Context, headers ...*common.Header) (response models.AccountApplicationResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s/applications/%s", common.EscapeParams(s.address, s.applicationId)...), s.p, headers)
	return
}
