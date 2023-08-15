package models

// SimulateUnnamedResourcesAccessed these are resources that were accessed by this
// group that would normally have caused failure, but were allowed in simulation.
// Depending on where this object is in the response, the unnamed resources it
// contains may or may not qualify for group resource sharing. If this is a field
// in SimulateTransactionGroupResult, the resources do qualify, but if this is a
// field in SimulateTransactionResult, they do not qualify. In order to make this
// group valid for actual submission, resources that qualify for group sharing can
// be made available by any transaction of the group; otherwise, resources must be
// placed in the same transaction which accessed them.
type SimulateUnnamedResourcesAccessed struct {
	// Accounts the unnamed accounts that were referenced. The order of this array is
	// arbitrary.
	Accounts []string `json:"accounts,omitempty"`

	// AppLocals the unnamed application local states that were referenced. The order
	// of this array is arbitrary.
	AppLocals []ApplicationLocalReference `json:"app-locals,omitempty"`

	// Apps the unnamed applications that were referenced. The order of this array is
	// arbitrary.
	Apps []uint64 `json:"apps,omitempty"`

	// AssetHoldings the unnamed asset holdings that were referenced. The order of this
	// array is arbitrary.
	AssetHoldings []AssetHoldingReference `json:"asset-holdings,omitempty"`

	// Assets the unnamed assets that were referenced. The order of this array is
	// arbitrary.
	Assets []uint64 `json:"assets,omitempty"`

	// Boxes the unnamed boxes that were referenced. The order of this array is
	// arbitrary.
	Boxes []BoxReference `json:"boxes,omitempty"`

	// ExtraBoxRefs the number of extra box references used to increase the IO budget.
	// This is in addition to the references defined in the input transaction group and
	// any referenced to unnamed boxes.
	ExtraBoxRefs uint64 `json:"extra-box-refs,omitempty"`
}
