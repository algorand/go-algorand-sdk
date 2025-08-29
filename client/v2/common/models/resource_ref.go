package models

// ResourceRef resourceRef names a single resource. Only one of the fields should
// be set.
type ResourceRef struct {
	// Address (d) Account whose balance record is accessible by the executing
	// ApprovalProgram or ClearStateProgram.
	Address string `json:"address,omitempty"`

	// ApplicationId (p) Application id whose GlobalState may be read by the executing
	// ApprovalProgram or ClearStateProgram.
	ApplicationId uint64 `json:"application-id,omitempty"`

	// AssetId (s) Asset whose AssetParams may be read by the executing
	// ApprovalProgram or ClearStateProgram.
	AssetId uint64 `json:"asset-id,omitempty"`

	// Box boxReference names a box by its name and the application ID it belongs to.
	Box BoxReference `json:"box,omitempty"`

	// Holding holdingRef names a holding by referring to an Address and Asset it
	// belongs to.
	Holding HoldingRef `json:"holding,omitempty"`

	// Local localsRef names a local state by referring to an Address and App it
	// belongs to.
	Local LocalsRef `json:"local,omitempty"`
}
