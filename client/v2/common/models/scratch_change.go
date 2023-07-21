package models

// ScratchChange a write operation into a scratch slot.
type ScratchChange struct {
	// NewValue represents an AVM value.
	NewValue AvmValue `json:"new-value"`

	// Slot the scratch slot written.
	Slot uint64 `json:"slot"`
}
