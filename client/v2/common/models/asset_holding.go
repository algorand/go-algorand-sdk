package models

// AssetHolding describes an asset held by an account.
// Definition:
// data/basics/userBalance.go : AssetHolding
type AssetHolding struct {
	// Amount (a) number of units held.
	Amount uint64 `json:"amount"`

	// AssetId asset ID of the holding.
	AssetId uint64 `json:"asset-id"`

	// IsFrozen (f) whether or not the holding is frozen.
	IsFrozen bool `json:"is-frozen"`
}
