package indexer

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

// LookupAssetByIDParams contains all of the query parameters for url serialization.
type LookupAssetByIDParams struct {

	// IncludeAll include all items including closed accounts, deleted applications,
	// destroyed assets, opted-out asset holdings, and closed-out application
	// localstates.
	IncludeAll bool `url:"include-all,omitempty"`
}

// LookupAssetByID lookup asset information.
type LookupAssetByID struct {
	c *Client

	assetId uint64

	p LookupAssetByIDParams
}

// IncludeAll include all items including closed accounts, deleted applications,
// destroyed assets, opted-out asset holdings, and closed-out application
// localstates.
func (s *LookupAssetByID) IncludeAll(IncludeAll bool) *LookupAssetByID {
	s.p.IncludeAll = IncludeAll

	return s
}

// Do performs the HTTP request
func (s *LookupAssetByID) Do(ctx context.Context, headers ...*common.Header) (validRound uint64, result models.Asset, err error) {
	response := models.AssetResponse{}
	err = s.c.get(ctx, &response, fmt.Sprintf("/v2/assets/%d", s.assetId), s.p, headers)
	validRound = response.CurrentRound
	result = response.Asset
	return
}
