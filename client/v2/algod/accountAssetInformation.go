package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// AccountAssetInformationParams contains all of the query parameters for url serialization.
type AccountAssetInformationParams struct {

	// Format configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

// AccountAssetInformation given a specific account public key, this call returns
// the account's asset holding and asset params (if they exist) about a given asset
// ID
type AccountAssetInformation struct {
	c *Client

	address string
	assetId uint64

	p AccountAssetInformationParams
}

// Do performs the HTTP request
func (s *AccountAssetInformation) Do(ctx context.Context, headers ...*common.Header) (response models.AccountAssetResponse, err error) {
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/accounts/%v/assets/%v", s.address, s.assetId), s.p, headers)
	return
}
