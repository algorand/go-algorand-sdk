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
	Creator string `json:"creator"`

	// Decimals (dc) The number of digits to use after the decimal point when
	// displaying this asset. If 0, the asset is not divisible. If 1, the base unit of
	// the asset is in tenths. If 2, the base unit of the asset is in hundredths, and
	// so on. This value must be between 0 and 19 (inclusive).
	Decimals uint64 `json:"decimals"`

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

	// Name (an) Name of this asset, as supplied by the creator. Included only when the
	// asset name is composed of printable utf-8 characters.
	Name string `json:"name,omitempty"`

	// NameB64 base64 encoded name of this asset, as supplied by the creator.
	NameB64 []byte `json:"name-b64,omitempty"`

	// Reserve (r) Address of account holding reserve (non-minted) units of this asset.
	Reserve string `json:"reserve,omitempty"`

	// Total (t) The total number of units of this asset.
	Total uint64 `json:"total"`

	// UnitName (un) Name of a unit of this asset, as supplied by the creator. Included
	// only when the name of a unit of this asset is composed of printable utf-8
	// characters.
	UnitName string `json:"unit-name,omitempty"`

	// UnitNameB64 base64 encoded name of a unit of this asset, as supplied by the
	// creator.
	UnitNameB64 []byte `json:"unit-name-b64,omitempty"`

	// Url (au) URL where more information about the asset can be retrieved. Included
	// only when the URL is composed of printable utf-8 characters.
	Url string `json:"url,omitempty"`

	// UrlB64 base64 encoded URL where more information about the asset can be
	// retrieved.
	UrlB64 []byte `json:"url-b64,omitempty"`
}
