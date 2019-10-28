package types

const assetUnitNameLen = 8
const assetNameLen = 32
const assetURLLen = 32
const assetMetadataHashLen = 32

// AssetIndex is the unique integer index of an asset that can be used to look
// up the creator of the asset, whose balance record contains the AssetParams
type AssetIndex uint64

// AssetParams describes the parameters of an asset.
type AssetParams struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Total specifies the total number of units of this asset
	// created.
	Total uint64 `codec:"t"`

	// DefaultFrozen specifies whether slots for this asset
	// in user accounts are frozen by default or not.
	DefaultFrozen bool `codec:"df"`

	// UnitName specifies a hint for the name of a unit of
	// this asset.
	UnitName [assetUnitNameLen]byte `codec:"un"`

	// AssetName specifies a hint for the name of the asset.
	AssetName [assetNameLen]byte `codec:"an"`

	// URL specifies a URL where more information about the asset can be
	// retrieved
	URL [assetURLLen]byte `codec:"au"`

	// MetadataHash specifies a commitment to some unspecified asset
	// metadata. The format of this metadata is up to the application.
	MetadataHash [assetMetadataHashLen]byte `codec:"am"`

	// Manager specifies an account that is allowed to change the
	// non-zero addresses in this AssetParams.
	Manager Address `codec:"m"`

	// Reserve specifies an account whose holdings of this asset
	// should be reported as "not minted".
	Reserve Address `codec:"r"`

	// Freeze specifies an account that is allowed to change the
	// frozen state of holdings of this asset.
	Freeze Address `codec:"f"`

	// Clawback specifies an account that is allowed to take units
	// of this asset from any account.
	Clawback Address `codec:"c"`
}
