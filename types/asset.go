package types

// AssetNameMaxLen is the max length in bytes for the asset name
const AssetNameMaxLen = 32

// AssetUnitNameMaxLen is the max length in bytes for the asset unit name
const AssetUnitNameMaxLen = 8

// AssetURLMaxLen is the max length in bytes for the asset url
const AssetURLMaxLen = 32

// AssetMetadataHashLen is the length of the AssetMetadataHash in bytes
const AssetMetadataHashLen = 32

// AssetMaxNumberOfDecimals is the maximum value of the Decimals field
const AssetMaxNumberOfDecimals = 19

// AssetIndex is the unique integer index of an asset that can be used to look
// up the creator of the asset, whose balance record contains the AssetParams
type AssetIndex uint64

// AssetParams describes the parameters of an asset.
type AssetParams struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Total specifies the total number of units of this asset
	// created.
	Total uint64 `codec:"t"`

	// Decimals specifies the number of digits to display after the decimal
	// place when displaying this asset. A value of 0 represents an asset
	// that is not divisible, a value of 1 represents an asset divisible
	// into tenths, and so on. This value must be between 0 and 19
	// (inclusive).
	Decimals uint32 `codec:"dc"`

	// DefaultFrozen specifies whether slots for this asset
	// in user accounts are frozen by default or not.
	DefaultFrozen bool `codec:"df"`

	// UnitName specifies a hint for the name of a unit of
	// this asset.
	UnitName string `codec:"un"`

	// AssetName specifies a hint for the name of the asset.
	AssetName string `codec:"an"`

	// URL specifies a URL where more information about the asset can be
	// retrieved
	URL string `codec:"au"`

	// MetadataHash specifies a commitment to some unspecified asset
	// metadata. The format of this metadata is up to the application.
	MetadataHash [AssetMetadataHashLen]byte `codec:"am"`

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
