package models

// AvmValue represents an AVM value.
type AvmValue struct {
	// Bytes bytes value.
	Bytes []byte `json:"bytes,omitempty"`

	// Type value type. Value `1` refers to **bytes**, value `2` refers to **uint64**
	Type uint64 `json:"type"`

	// Uint uint value.
	Uint uint64 `json:"uint,omitempty"`
}
