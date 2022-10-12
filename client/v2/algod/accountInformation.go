package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// AccountInformationParams contains all of the query parameters for url serialization.
type AccountInformationParams struct {

	// Exclude when set to `all` will exclude asset holdings, application local state,
	// created asset parameters, any created application parameters. Defaults to
	// `none`.
	Exclude string `url:"exclude,omitempty"`

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

// AccountInformation given a specific account public key, this call returns the
// accounts status, balance and spendable amounts
type AccountInformation struct {
	c *Client

	address string

	p AccountInformationParams
}

// Exclude when set to `all` will exclude asset holdings, application local state,
// created asset parameters, any created application parameters. Defaults to
// `none`.
func (s *AccountInformation) Exclude(Exclude string) *AccountInformation {
	s.p.Exclude = Exclude

	return s
}

// Do performs the HTTP request
func (s *AccountInformation) Do(ctx context.Context, headers ...*common.Header) (response models.Account, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s", common.EscapeParams(s.address)...), s.p, headers)
	return
}
