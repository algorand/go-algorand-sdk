package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// AccountInformationParams contains all of the query parameters for url serialization.
type AccountInformationParams struct {

	// Exclude exclude additional items from the account. Use `all` to exclude asset
	// holdings, application local state, created asset parameters, and created
	// application parameters. Use `created-apps-params` to exclude only the parameters
	// of created applications (returns only application IDs). Use
	// `created-assets-params` to exclude only the parameters of created assets
	// (returns only asset IDs). Multiple values can be comma-separated (e.g.,
	// `created-apps-params,created-assets-params`). Note: `all` and `none` cannot be
	// combined with other values. Defaults to `none`.
	Exclude []string `url:"exclude,omitempty,comma"`

	// Format configures whether the response object is JSON or MessagePack encoded. If
	// not provided, defaults to JSON.
	Format string `url:"format,omitempty"`
}

// AccountInformation given a specific account public key, this call returns the
// account's status, balance and spendable amounts
type AccountInformation struct {
	c *Client

	address string

	p AccountInformationParams
}

// Exclude exclude additional items from the account. Use `all` to exclude asset
// holdings, application local state, created asset parameters, and created
// application parameters. Use `created-apps-params` to exclude only the parameters
// of created applications (returns only application IDs). Use
// `created-assets-params` to exclude only the parameters of created assets
// (returns only asset IDs). Multiple values can be comma-separated (e.g.,
// `created-apps-params,created-assets-params`). Note: `all` and `none` cannot be
// combined with other values. Defaults to `none`.
func (s *AccountInformation) Exclude(Exclude []string) *AccountInformation {
	s.p.Exclude = Exclude

	return s
}

// Do performs the HTTP request
func (s *AccountInformation) Do(ctx context.Context, headers ...*common.Header) (response models.Account, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%s", common.EscapeParams(s.address)...), s.p, headers)
	return
}
