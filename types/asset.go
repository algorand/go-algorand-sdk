package types

const assetUnitNameLen = 8
const assetNameLen = 32

// AssetID is a name of an asset.
type AssetID struct {
	Creator Address `codec:"c"`
	Index   uint64  `codec:"i"`
}

// AssetParams describes the parameters of an asset.
type AssetParams struct {
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
