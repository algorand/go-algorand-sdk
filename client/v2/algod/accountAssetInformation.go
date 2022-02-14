package algod

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// AccountAssetInformationParams defines parameters for AccountAssetInformation.
type AccountAssetInformationParams struct {
	// Configures whether the response object is JSON or MessagePack encoded.
	Format string `url:"format,omitempty"`
}

type AccountAssetInformation struct {
	c *Client

	address string
	assetID uint64

	p AccountAssetInformationParams
}

// Do performs the HTTP request.
func (s *AccountAssetInformation) Do(ctx context.Context, headers ...*common.Header) (response models.AccountAssetResponse, err error) {
	err = s.c.get(
		ctx, &response, fmt.Sprintf("/v2/accounts/%s/assets/%d", s.address, s.assetID),
		s.p, headers)
	return
}
