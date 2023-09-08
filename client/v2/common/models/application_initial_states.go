package models

// ApplicationInitialStates an application's initial global/local/box states that
// were accessed during simulation.
type ApplicationInitialStates struct {
	// AppBoxes an application's global/local/box state.
	AppBoxes ApplicationKVStorage `json:"app-boxes,omitempty"`

	// AppGlobals an application's global/local/box state.
	AppGlobals ApplicationKVStorage `json:"app-globals,omitempty"`

	// AppLocals an application's initial local states tied to different accounts.
	AppLocals []ApplicationKVStorage `json:"app-locals,omitempty"`

	// Id application index.
	Id uint64 `json:"id"`
}
