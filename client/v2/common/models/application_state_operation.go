package models

// ApplicationStateOperation an operation against an application's global/local/box
// state.
type ApplicationStateOperation struct {
	// Account for local state changes, the address of the account associated with the
	// local state.
	Account string `json:"account,omitempty"`

	// AppStateType type of application state. Value `g` is **global state**, `l` is
	// **local state**, `b` is **boxes**.
	AppStateType string `json:"app-state-type"`

	// Key the key (name) of the global/local/box state.
	Key []byte `json:"key"`

	// NewValue represents an AVM value.
	NewValue AvmValue `json:"new-value,omitempty"`

	// Operation operation type. Value `w` is **write**, `d` is **delete**.
	Operation string `json:"operation"`
}
