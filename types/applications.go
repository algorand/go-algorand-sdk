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

// BoxReference names a box by the index in the foreign app array
type BoxReference struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// The index of the app in the foreign app array.
	ForeignAppIdx uint64 `codec:"i"`

	// The name of the box unique to the app it belongs to
	Name []byte `codec:"n"`
}

const (
	// encodedMaxApplicationArgs sets the allocation bound for the maximum
	// number of ApplicationArgs that a transaction decoded off of the wire
	// can contain. Its value is verified against consensus parameters in
	// TestEncodedAppTxnAllocationBounds
	encodedMaxApplicationArgs = 32

	// encodedMaxAccounts sets the allocation bound for the maximum number
	// of Accounts that a transaction decoded off of the wire can contain.
	// Its value is verified against consensus parameters in
	// TestEncodedAppTxnAllocationBounds
	encodedMaxAccounts = 32

	// encodedMaxForeignApps sets the allocation bound for the maximum
	// number of ForeignApps that a transaction decoded off of the wire can
	// contain. Its value is verified against consensus parameters in
	// TestEncodedAppTxnAllocationBounds
	encodedMaxForeignApps = 32

	// encodedMaxForeignAssets sets the allocation bound for the maximum
	// number of ForeignAssets that a transaction decoded off of the wire
	// can contain. Its value is verified against consensus parameters in
	// TestEncodedAppTxnAllocationBounds
	encodedMaxForeignAssets = 32

	// encodedMaxBoxReferences sets the allocation bound for the maximum
	// number of BoxReferences that a transaction decoded off of the wire
	// can contain. Its value is verified against consensus parameters in
	// TestEncodedAppTxnAllocationBounds
	encodedMaxBoxReferences = 32
)

// OnCompletion is an enum representing some layer 1 side effect that an
// ApplicationCall transaction will have if it is included in a block.
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

	ApplicationID   AppIndex       `codec:"apid"`
	OnCompletion    OnCompletion   `codec:"apan"`
	ApplicationArgs [][]byte       `codec:"apaa,allocbound=encodedMaxApplicationArgs"`
	Accounts        []Address      `codec:"apat,allocbound=encodedMaxAccounts"`
	ForeignApps     []AppIndex     `codec:"apfa,allocbound=encodedMaxForeignApps"`
	ForeignAssets   []AssetIndex   `codec:"apas,allocbound=encodedMaxForeignAssets"`
	BoxReferences   []BoxReference `codec:"apbx,allocbound=encodedMaxBoxReferences"`

	LocalStateSchema  StateSchema `codec:"apls"`
	GlobalStateSchema StateSchema `codec:"apgs"`
	ApprovalProgram   []byte      `codec:"apap"`
	ClearStateProgram []byte      `codec:"apsu"`
	ExtraProgramPages uint32      `codec:"apep"`

	// If you add any fields here, remember you MUST modify the Empty
	// method below!
}

// StateSchema sets maximums on the number of each type that may be stored
type StateSchema struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	NumUint      uint64 `codec:"nui"`
	NumByteSlice uint64 `codec:"nbs"`
}

// Empty indicates whether or not all the fields in the
// ApplicationCallTxnFields are zeroed out
func (ac *ApplicationCallTxnFields) Empty() bool {
	if ac.ApplicationID != 0 {
		return false
	}
	if ac.OnCompletion != 0 {
		return false
	}
	if ac.ApplicationArgs != nil {
		return false
	}
	if ac.Accounts != nil {
		return false
	}
	if ac.ForeignApps != nil {
		return false
	}
	if ac.ForeignAssets != nil {
		return false
	}
	if ac.BoxReferences != nil {
		return false
	}
	if ac.LocalStateSchema != (StateSchema{}) {
		return false
	}
	if ac.GlobalStateSchema != (StateSchema{}) {
		return false
	}
	if ac.ApprovalProgram != nil {
		return false
	}
	if ac.ClearStateProgram != nil {
		return false
	}
	if ac.ExtraProgramPages != 0 {
		return false
	}
	return true
}
