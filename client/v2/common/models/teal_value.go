package models

// TealValue represents a TEAL value.
type TealValue struct {
	// Bytes (tb) bytes value.
	Bytes string `json:"bytes,omitempty"`

	// Type (tt) value type.
	Type uint64 `json:"type,omitempty"`

	// Uint (ui) uint value.
	Uint uint64 `json:"uint,omitempty"`
}
