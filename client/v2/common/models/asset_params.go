package models

// AssetParams assetParams specifies the parameters for an asset.
// (apar) when part of an AssetConfig transaction.
// Definition:
// data/transactions/asset.go : AssetParams
type AssetParams struct {
	// Clawback (c) Address of account used to clawback holdings of this asset. If
	// empty, clawback is not permitted.
	Clawback string `json:"clawback,omitempty"`

	// Creator the address that created this asset. This is the address where the
	// parameters for this asset can be found, and also the address where unwanted
	// asset units can be sent in the worst case.
	Creator string `json:"creator,omitempty"`

	// Decimals (dc) The number of digits to use after the decimal point when
	// displaying this asset. If 0, the asset is not divisible. If 1, the base unit of
	// the asset is in tenths. If 2, the base unit of the asset is in hundredths, and
	// so on. This value must be between 0 and 19 (inclusive).
	Decimals uint64 `json:"decimals,omitempty"`

	// DefaultFrozen (df) Whether holdings of this asset are frozen by default.
	DefaultFrozen bool `json:"default-frozen,omitempty"`

	// Freeze (f) Address of account used to freeze holdings of this asset. If empty,
	// freezing is not permitted.
	Freeze string `json:"freeze,omitempty"`

	// Manager (m) Address of account used to manage the keys of this asset and to
	// destroy it.
	Manager string `json:"manager,omitempty"`

	// MetadataHash (am) A commitment to some unspecified asset metadata. The format of
	// this metadata is up to the application.
	MetadataHash []byte `json:"metadata-hash,omitempty"`

	// Name (an) Name of this asset, as supplied by the creator.
	Name string `json:"name,omitempty"`

	// Reserve (r) Address of account holding reserve (non-minted) units of this asset.
	Reserve string `json:"reserve,omitempty"`

	// Total (t) The total number of units of this asset.
	Total uint64 `json:"total,omitempty"`

	// UnitName (un) Name of a unit of this asset, as supplied by the creator.
	UnitName string `json:"unit-name,omitempty"`

	// Url (au) URL where more information about the asset can be retrieved.
	Url string `json:"url,omitempty"`
}
