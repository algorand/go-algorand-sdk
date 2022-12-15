package models

// AssetResourceRecord represents AssetParams and AssetHolding in deltas
type AssetResourceRecord struct {
	// Address account address of the asset
	Address string `json:"address"`

	// AssetDeleted whether the asset was deleted
	AssetDeleted bool `json:"asset-deleted"`

	// AssetHolding the asset holding
	AssetHolding AssetHolding `json:"asset-holding,omitempty"`

	// AssetHoldingDeleted whether the asset holding was deleted
	AssetHoldingDeleted bool `json:"asset-holding-deleted"`

	// AssetIndex index of the asset
	AssetIndex uint64 `json:"asset-index"`

	// AssetParams asset params
	AssetParams AssetParams `json:"asset-params,omitempty"`
}
