package models

// Asset specifies both the unique identifier and the parameters for an asset
type Asset struct {
	// Index unique asset identifier
	Index uint64 `json:"index"`

	// Params assetParams specifies the parameters for an asset.
	// (apar) when part of an AssetConfig transaction.
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params AssetParams `json:"params"`
}
