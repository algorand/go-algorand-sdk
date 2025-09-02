package types

// This file has the applications specific structures

// ApplicationFields are the fields that are common to all application
type ApplicationFields struct {
	ApplicationCallTxnFields
}

// AppIndex is the unique integer index of an application that can be used to
// look up the creator of the application, whose balance record contains the
// AppParams
type AppIndex uint64

// AppBoxReference names a box by the app ID
type AppBoxReference struct {
	// The ID of the app that owns the box. Must be converted to BoxReference during transaction submission.
	AppID uint64

	// The Name of the box unique to the app it belongs to
	Name []byte
}

// AppHoldingRef identifies an asset holding by the asset id and the address (zero address/empty string means Sender).
// It can be viewed as the "hydrated" form of a HoldingRef, which uses indices.
type AppHoldingRef struct {
	Asset   uint64
	Address string // empty string means Sender
}

// AppLocalsRef identifies local state by the app id and the address (zero address/empty string means Sender).
// It can be viewed as the "hydrated" form of a LocalsRef, which uses indices.
type AppLocalsRef struct {
	App     uint64
	Address string // empty string means Sender
}

// BoxReference names a box by the index in the foreign app array
type BoxReference struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// The index of the app in the foreign app array.
	ForeignAppIdx uint64 `codec:"i"`

	// The name of the box unique to the app it belongs to
	Name []byte `codec:"n"`
}

// ResourceRef names a single resource
type ResourceRef struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Only one of these may be set
	Address Address      `codec:"d"`
	Asset   AssetIndex   `codec:"s"`
	App     AppIndex     `codec:"p"`
	Holding HoldingRef   `codec:"h"`
	Locals  LocalsRef    `codec:"l"`
	Box     BoxReference `codec:"b"`
}

// HoldingRef names a holding by referring to an Address and Asset that appear
// earlier in the Access list (0 is special cased)
type HoldingRef struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
	Address uint64   `codec:"d"` // 0=Sender,n-1=index into the Access list, which must be an Address
	Asset   uint64   `codec:"s"` // n-1=index into the Access list, which must be an Asset
}

// LocalsRef names a local state by referring to an Address and App that appear
// earlier in the Access list (0 is special cased)
type LocalsRef struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
	Address uint64   `codec:"d"` // 0=Sender,n-1=index into the Access list, which must be an Address
	App     uint64   `codec:"p"` // 0=ApplicationID,n-1=index into the Access list, which must be an App
}

// OnCompletion is an enum representing some layer 1 side effect that an
// ApplicationCall transaction will have if it is included in a block.
//
//go:generate stringer -type=OnCompletion -output=application_string.go
type OnCompletion uint64

const (
	// NoOpOC indicates that an application transaction will simply call its
	// ApprovalProgram
	NoOpOC OnCompletion = 0

	// OptInOC indicates that an application transaction will allocate some
	// LocalState for the application in the sender's account
	OptInOC OnCompletion = 1

	// CloseOutOC indicates that an application transaction will deallocate
	// some LocalState for the application from the user's account
	CloseOutOC OnCompletion = 2

	// ClearStateOC is similar to CloseOutOC, but may never fail. This
	// allows users to reclaim their minimum balance from an application
	// they no longer wish to opt in to.
	ClearStateOC OnCompletion = 3

	// UpdateApplicationOC indicates that an application transaction will
	// update the ApprovalProgram and ClearStateProgram for the application
	UpdateApplicationOC OnCompletion = 4

	// DeleteApplicationOC indicates that an application transaction will
	// delete the AppParams for the application from the creator's balance
	// record
	DeleteApplicationOC OnCompletion = 5
)

// ApplicationCallTxnFields captures the transaction fields used for all
// interactions with applications
type ApplicationCallTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// ApplicationID is 0 when creating an application, and nonzero when
	// calling an existing application.
	ApplicationID AppIndex `codec:"apid"`

	// OnCompletion specifies an optional side-effect that this transaction
	// will have on the balance record of the sender or the application's
	// creator. See the documentation for the OnCompletion type for more
	// information on each possible value.
	OnCompletion OnCompletion `codec:"apan"`

	// ApplicationArgs are arguments accessible to the executing
	// ApprovalProgram or ClearStateProgram.
	ApplicationArgs [][]byte `codec:"apaa"`

	// Accounts are accounts whose balance records are accessible
	// by the executing ApprovalProgram or ClearStateProgram. To
	// access LocalState or an ASA balance for an account besides
	// the sender, that account's address must be listed here (and
	// since v4, the ForeignApp or ForeignAsset must also include
	// the app or asset id).
	Accounts []Address `codec:"apat"`

	// ForeignAssets are asset IDs for assets whose AssetParams
	// (and since v4, Holdings) may be read by the executing
	// ApprovalProgram or ClearStateProgram.
	ForeignAssets []AssetIndex `codec:"apas"`

	// ForeignApps are application IDs for applications besides
	// this one whose GlobalState (or Local, since v4) may be read
	// by the executing ApprovalProgram or ClearStateProgram.
	ForeignApps []AppIndex `codec:"apfa"`

	// Access unifies `Accounts`, `ForeignApps`, `ForeignAssets`, and `Boxes`
	// under a single list. It removes all implicitly available resources, so
	// "cross-product" resources (holdings and locals) must be explicitly
	// listed, as well as app accounts, even the app account of the called app!
	// Transactions using Access may not use the other lists.
	Access []ResourceRef `codec:"al"`

	// Boxes are the boxes that can be accessed by this transaction (and others
	// in the same group). The Index in the BoxRef is the slot of ForeignApps
	// that the name is associated with (shifted by 1. 0 indicates "current
	// app")
	BoxReferences []BoxReference `codec:"apbx"`

	// LocalStateSchema specifies the maximum number of each type that may
	// appear in the local key/value store of users who opt in to this
	// application. This field is only used during application creation
	// (when the ApplicationID field is 0),
	LocalStateSchema StateSchema `codec:"apls"`

	// GlobalStateSchema specifies the maximum number of each type that may
	// appear in the global key/value store associated with this
	// application. This field is only used during application creation
	// (when the ApplicationID field is 0).
	GlobalStateSchema StateSchema `codec:"apgs"`

	// ApprovalProgram is the stateful TEAL bytecode that executes on all
	// ApplicationCall transactions associated with this application,
	// except for those where OnCompletion is equal to ClearStateOC. If
	// this program fails, the transaction is rejected. This program may
	// read and write local and global state for this application.
	ApprovalProgram []byte `codec:"apap"`

	// ClearStateProgram is the stateful TEAL bytecode that executes on
	// ApplicationCall transactions associated with this application when
	// OnCompletion is equal to ClearStateOC. This program will not cause
	// the transaction to be rejected, even if it fails. This program may
	// read and write local and global state for this application.
	ClearStateProgram []byte `codec:"apsu"`

	// ExtraProgramPages specifies the additional app program len requested in pages.
	// A page is MaxAppProgramLen bytes. This field enables execution of app programs
	// larger than the default config, MaxAppProgramLen.
	ExtraProgramPages uint32 `codec:"apep"`

	// RejectVersion is the lowest application version for which this
	// transaction should immediately fail. 0 indicates that no version check should be performed.
	RejectVersion uint64 `codec:"aprv"`

	// If you add any fields here, remember you MUST modify the Empty
	// method below!
}

// StateSchema sets maximums on the number of each type that may be stored
type StateSchema struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	NumUint      uint64 `codec:"nui"`
	NumByteSlice uint64 `codec:"nbs"`
}
