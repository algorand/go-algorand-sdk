package models

// TealValue represents a TEAL value.
type TealValue struct {
	// Bytes (tb) bytes value.
	Bytes string `json:"bytes"`

	// Type (tt) value type. Value `1` refers to **bytes**, value `2` refers to
	// **uint**
	Type uint64 `json:"type"`

	// Uint (ui) uint value.
	Uint uint64 `json:"uint"`
}
