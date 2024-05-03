package models

// AccountAssetHolding accountAssetHolding describes the account's asset holding
// and asset parameters (if either exist) for a specific asset ID.
type AccountAssetHolding struct {
	// AssetHolding (asset) Details about the asset held by this account.
	// The raw account uses `AssetHolding` for this type.
	AssetHolding AssetHolding `json:"asset-holding"`

	// AssetParams (apar) parameters of the asset held by this account.
	// The raw account uses `AssetParams` for this type.
	AssetParams AssetParams `json:"asset-params,omitempty"`
}
