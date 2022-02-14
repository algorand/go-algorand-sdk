package models

// AccountAssetResponse defines model for AccountAssetResponse.
type AccountAssetResponse struct {
	// Describes an asset held by an account.
	//
	// Definition:
	// data/basics/userBalance.go : AssetHolding
	AssetHolding *AssetHolding `json:"asset-holding,omitempty"`

	// AssetParams specifies the parameters for an asset.
	//
	// \[apar\] when part of an AssetConfig transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetParams
	CreatedAsset *AssetParams `json:"created-asset,omitempty"`

	// The round for which this information is relevant.
	Round uint64 `json:"round"`
}
