package models

// MiniAssetHolding a simplified version of AssetHolding
type MiniAssetHolding struct {
	// Address
	Address string `json:"address"`

	// Amount
	Amount uint64 `json:"amount"`

	// Deleted whether or not this asset holding is currently deleted from its account.
	Deleted bool `json:"deleted,omitempty"`

	// IsFrozen
	IsFrozen bool `json:"is-frozen"`

	// OptedInAtRound round during which the account opted into the asset.
	OptedInAtRound uint64 `json:"opted-in-at-round,omitempty"`

	// OptedOutAtRound round during which the account opted out of the asset.
	OptedOutAtRound uint64 `json:"opted-out-at-round,omitempty"`
}
