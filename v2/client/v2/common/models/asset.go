package models

// Asset specifies both the unique identifier and the parameters for an asset
type Asset struct {
	// CreatedAtRound round during which this asset was created.
	CreatedAtRound uint64 `json:"created-at-round,omitempty"`

	// Deleted whether or not this asset is currently deleted.
	Deleted bool `json:"deleted,omitempty"`

	// DestroyedAtRound round during which this asset was destroyed.
	DestroyedAtRound uint64 `json:"destroyed-at-round,omitempty"`

	// Index unique asset identifier
	Index uint64 `json:"index"`

	// Params assetParams specifies the parameters for an asset.
	// (apar) when part of an AssetConfig transaction.
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params AssetParams `json:"params"`
}
