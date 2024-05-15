package models

// AccountAssetsInformationResponse accountAssetsInformationResponse contains a
// list of assets held by an account.
type AccountAssetsInformationResponse struct {
	// AssetHoldings
	AssetHoldings []AccountAssetHolding `json:"asset-holdings,omitempty"`

	// NextToken used for pagination, when making another request provide this token
	// with the next parameter.
	NextToken string `json:"next-token,omitempty"`

	// Round the round for which this information is relevant.
	Round uint64 `json:"round"`
}
