package models

// ApplicationParams stores the global information associated with an application.
type ApplicationParams struct {
	// ApprovalProgram (approv) approval program.
	ApprovalProgram []byte `json:"approval-program"`

	// ClearStateProgram (clearp) approval program.
	ClearStateProgram []byte `json:"clear-state-program"`

	// Creator the address that created this application. This is the address where the
	// parameters and global state for this application can be found.
	Creator string `json:"creator,omitempty"`

	// ExtraProgramPages (epp) the amount of extra program pages available to this app.
	ExtraProgramPages uint64 `json:"extra-program-pages,omitempty"`

	// GlobalState [\gs) global schema
	GlobalState []TealKeyValue `json:"global-state,omitempty"`

	// GlobalStateSchema [\gsch) global schema
	GlobalStateSchema ApplicationStateSchema `json:"global-state-schema,omitempty"`

	// LocalStateSchema [\lsch) local schema
	LocalStateSchema ApplicationStateSchema `json:"local-state-schema,omitempty"`
}
