package models

// TealValue represents a TEAL value.
type TealValue struct {
	// Bytes bytes value.
	Bytes string `json:"bytes"`

	// Type type of the value. Value `1` refers to **bytes**, value `2` refers to
	// **uint**
	Type uint64 `json:"type"`

	// Uint uint value.
	Uint uint64 `json:"uint"`
}
