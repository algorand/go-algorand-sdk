package models

// TransactionAssetFreeze fields for an asset freeze transaction.
// Definition:
// data/transactions/asset.go : AssetFreezeTxnFields
type TransactionAssetFreeze struct {
	// Address (fadd) Address of the account whose asset is being frozen or thawed.
	Address string `json:"address"`

	// AssetId (faid) ID of the asset being frozen or thawed.
	AssetId uint64 `json:"asset-id"`

	// NewFreezeStatus (afrz) The new freeze status.
	NewFreezeStatus bool `json:"new-freeze-status"`
}
