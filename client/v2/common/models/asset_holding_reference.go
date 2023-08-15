package models

// AssetHoldingReference references an asset held by an account.
type AssetHoldingReference struct {
	// Account address of the account holding the asset.
	Account string `json:"account"`

	// Asset asset ID of the holding.
	Asset uint64 `json:"asset"`
}
