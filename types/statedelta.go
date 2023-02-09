package types

// TealType is an enum of the types in a TEAL program: Bytes and Uint
type TealType uint64

// TealValue contains type information and a value, representing a value in a
// TEAL program
type TealValue struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Type  TealType `codec:"tt"`
	Bytes string   `codec:"tb"`
	Uint  uint64   `codec:"ui"`
}

// TealKeyValue represents a key/value store for use in an application's
// LocalState or GlobalState
//msgp:allocbound TealKeyValue EncodedMaxKeyValueEntries
type TealKeyValue map[string]TealValue

// StateSchemas is a thin wrapper around the LocalStateSchema and the
// GlobalStateSchema, since they are often needed together
type StateSchemas struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	LocalStateSchema  StateSchema `codec:"lsch"`
	GlobalStateSchema StateSchema `codec:"gsch"`
}

// AppParams stores the global information associated with an application
type AppParams struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	ApprovalProgram   []byte       `codec:"approv,allocbound=config.MaxAvailableAppProgramLen"`
	ClearStateProgram []byte       `codec:"clearp,allocbound=config.MaxAvailableAppProgramLen"`
	GlobalState       TealKeyValue `codec:"gs"`
	StateSchemas
	ExtraProgramPages uint32 `codec:"epp"`
}

// AppLocalState stores the LocalState associated with an application. It also
// stores a cached copy of the application's LocalStateSchema so that
// MinBalance requirements may be computed 1. without looking up the
// AppParams and 2. even if the application has been deleted
type AppLocalState struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Schema   StateSchema  `codec:"hsch"`
	KeyValue TealKeyValue `codec:"tkv"`
}

// AppLocalStateDelta tracks a changed AppLocalState, and whether it was deleted
type AppLocalStateDelta struct {
	LocalState *AppLocalState
	Deleted    bool
}

// AppParamsDelta tracks a changed AppParams, and whether it was deleted
type AppParamsDelta struct {
	Params  *AppParams
	Deleted bool
}

// AppResourceRecord represents AppParams and AppLocalState in deltas
type AppResourceRecord struct {
	Aidx   AppIndex
	Addr   Address
	Params AppParamsDelta
	State  AppLocalStateDelta
}

// AssetHolding describes an asset held by an account.
type AssetHolding struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Amount uint64 `codec:"a"`
	Frozen bool   `codec:"f"`
}

// AssetHoldingDelta records a changed AssetHolding, and whether it was deleted
type AssetHoldingDelta struct {
	Holding *AssetHolding
	Deleted bool
}

// AssetParamsDelta tracks a changed AssetParams, and whether it was deleted
type AssetParamsDelta struct {
	Params  *AssetParams
	Deleted bool
}

// AssetResourceRecord represents AssetParams and AssetHolding in deltas
type AssetResourceRecord struct {
	Aidx    AssetIndex
	Addr    Address
	Params  AssetParamsDelta
	Holding AssetHoldingDelta
}

// AccountAsset is used as a map key.
type AccountAsset struct {
	Address Address
	Asset   AssetIndex
}

// AccountApp is used as a map key.
type AccountApp struct {
	Address Address
	App     AppIndex
}

// A OneTimeSignatureVerifier is used to identify the holder of
// OneTimeSignatureSecrets and prove the authenticity of OneTimeSignatures
// against some OneTimeSignatureIdentifier.
type OneTimeSignatureVerifier [32]byte

// VRFVerifier is a deprecated name for VrfPubkey
type VRFVerifier [32]byte

// VotingData holds participation information
type VotingData struct {
	VoteID       OneTimeSignatureVerifier
	SelectionID  VRFVerifier
	StateProofID Commitment

	VoteFirstValid  Round
	VoteLastValid   Round
	VoteKeyDilution uint64
}

// Status is the delegation status of an account's MicroAlgos
type Status byte

// AccountBaseData contains base account info like balance, status and total number of resources
type AccountBaseData struct {
	Status             Status
	MicroAlgos         MicroAlgos
	RewardsBase        uint64
	RewardedMicroAlgos MicroAlgos
	AuthAddr           Address

	TotalAppSchema      StateSchema // Totals across created globals, and opted in locals.
	TotalExtraAppPages  uint32      // Total number of extra pages across all created apps
	TotalAppParams      uint64      // Total number of apps this account has created
	TotalAppLocalStates uint64      // Total number of apps this account is opted into.
	TotalAssetParams    uint64      // Total number of assets created by this account
	TotalAssets         uint64      // Total of asset creations and optins (i.e. number of holdings)
	TotalBoxes          uint64      // Total number of boxes associated to this account
	TotalBoxBytes       uint64      // Total bytes for this account's boxes. keys _and_ values count
}

// AccountData provides users of the Balances interface per-account data (like basics.AccountData)
// but without any maps containing AppParams, AppLocalState, AssetHolding, or AssetParams. This
// ensures that transaction evaluation must retrieve and mutate account, asset, and application data
// separately, to better support on-disk and in-memory schemas that do not store them together.
type AccountData struct {
	AccountBaseData
	VotingData
}

// BalanceRecord is similar to basics.BalanceRecord but with decoupled base and voting data
type BalanceRecord struct {
	Addr Address
	AccountData
}

// AccountDeltas stores ordered accounts and allows fast lookup by address
// One key design aspect here was to ensure that we're able to access the written
// deltas in a deterministic order, while maintaining O(1) lookup. In order to
// do that, each of the arrays here is constructed as a pair of (slice, map).
// The map would point the address/address+creatable id onto the index of the
// element within the slice.
// If adding fields make sure to add them to the .reset() method to avoid dirty state
type AccountDeltas struct {
	// Actual data. If an account is deleted, `Accts` contains the BalanceRecord
	// with an empty `AccountData` and a populated `Addr`.
	Accts []BalanceRecord
	// cache for addr to deltas index resolution
	acctsCache map[Address]int

	// AppResources deltas. If app params or local state is deleted, there is a nil value in AppResources.Params or AppResources.State and Deleted flag set
	AppResources []AppResourceRecord
	// caches for {addr, app id} to app params delta resolution
	// not preallocated - use UpsertAppResource instead of inserting directly
	appResourcesCache map[AccountApp]int

	AssetResources []AssetResourceRecord
	// not preallocated - use UpsertAssertResource instead of inserting directly
	assetResourcesCache map[AccountAsset]int
}

// A KvValueDelta shows how the Data associated with a key in the kvstore has
// changed.  However, OldData is elided during evaluation, and only filled in at
// the conclusion of a block during the called to roundCowState.deltas()
type KvValueDelta struct {
	// Data stores the most recent value (nil == deleted)
	Data []byte

	// OldData stores the previous vlaue (nil == didn't exist)
	OldData []byte
}

// IncludedTransactions defines the transactions included in a block, their index and last valid round.
type IncludedTransactions struct {
	LastValid Round
	Intra     uint64 // the index of the transaction in the block
}

// Txid is a hash used to uniquely identify individual transactions
type Txid Digest

// A Txlease is a transaction (sender, lease) pair which uniquely specifies a
// transaction lease.
type Txlease struct {
	Sender Address
	Lease  [32]byte
}

// CreatableIndex represents either an AssetIndex or AppIndex, which come from
// the same namespace of indices as each other (both assets and apps are
// "creatables")
type CreatableIndex uint64

// CreatableType is an enum representing whether or not a given creatable is an
// application or an asset
type CreatableType uint64

// ModifiedCreatable defines the changes to a single single creatable state
type ModifiedCreatable struct {
	// Type of the creatable: app or asset
	Ctype CreatableType

	// Created if true, deleted if false
	Created bool

	// creator of the app/asset
	Creator Address

	// Keeps track of how many times this app/asset appears in
	// accountUpdates.creatableDeltas
	Ndeltas int
}

// AlgoCount represents a total of algos of a certain class
// of accounts (split up by their Status value).
type AlgoCount struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Sum of algos of all accounts in this class.
	Money MicroAlgos `codec:"mon"`

	// Total number of whole reward units in accounts.
	RewardUnits uint64 `codec:"rwd"`
}

// AccountTotals represents the totals of algos in the system
// grouped by different account status values.
type AccountTotals struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Online           AlgoCount `codec:"online"`
	Offline          AlgoCount `codec:"offline"`
	NotParticipating AlgoCount `codec:"notpart"`

	// Total number of algos received per reward unit since genesis
	RewardsLevel uint64 `codec:"rwdlvl"`
}

// LedgerStateDelta describes the delta between a given round to the previous round
// If adding a new field not explicitly allocated by PopulateStateDelta, make sure to reset
// it in .ReuseStateDelta to avoid dirty memory errors.
// If adding fields make sure to add them to the .Reset() method to avoid dirty state
type LedgerStateDelta struct {
	// modified new accounts
	Accts AccountDeltas

	// modified kv pairs (nil == delete)
	// not preallocated use .AddKvMod to insert instead of direct assignment
	KvMods map[string]KvValueDelta

	// new Txids for the txtail and TxnCounter, mapped to txn.LastValid
	Txids map[Txid]IncludedTransactions

	// new txleases for the txtail mapped to expiration
	// not pre-allocated so use .AddTxLease to insert instead of direct assignment
	Txleases map[Txlease]Round

	// new creatables creator lookup table
	// not pre-allocated so use .AddCreatable to insert instead of direct assignment
	Creatables map[CreatableIndex]ModifiedCreatable

	// new block header; read-only
	Hdr *BlockHeader

	// next round for which we expect a state proof.
	// zero if no state proof is expected.
	StateProofNext Round

	// previous block timestamp
	PrevTimestamp int64

	// initial hint for allocating data structures for StateDelta
	initialHint int

	// The account totals reflecting the changes in this StateDelta object.
	Totals AccountTotals
}
