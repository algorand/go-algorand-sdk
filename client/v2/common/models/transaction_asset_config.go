package models

// TransactionAssetConfig fields for asset allocation, re-configuration, and
// destruction.
// A zero value for asset-id indicates asset creation.
// A zero value for the params indicates asset destruction.
// Definition:
// data/transactions/asset.go : AssetConfigTxnFields
type TransactionAssetConfig struct {
	// AssetId (xaid) ID of the asset being configured or empty if creating.
	AssetId uint64 `json:"asset-id,omitempty"`

	// Params assetParams specifies the parameters for an asset.
	// (apar) when part of an AssetConfig transaction.
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params AssetParams `json:"params,omitempty"`
}
