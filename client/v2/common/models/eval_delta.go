package models

// EvalDelta represents a TEAL value delta.
type EvalDelta struct {
	// Action (at) delta action.
	Action uint64 `json:"action"`

	// Bytes (bs) bytes value.
	Bytes string `json:"bytes,omitempty"`

	// Uint (ui) uint value.
	Uint uint64 `json:"uint,omitempty"`
}
