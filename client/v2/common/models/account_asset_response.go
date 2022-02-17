package models

// AccountAssetResponse accountAssetResponse describes the asset data for a
// specific account and asset ID.
type AccountAssetResponse struct {
	// AssetHolding (asset) Details about the asset held by this account.
	// The raw account uses `AssetHolding` for this type.
	AssetHolding AssetHolding `json:"asset-holding,omitempty"`

	// CreatedAsset (apar) parameters of the asset created by this account.
	// The raw account uses `AssetParams` for this type.
	CreatedAsset AssetParams `json:"created-asset,omitempty"`

	// Round the round for which this information is relevant.
	Round uint64 `json:"round"`
}
