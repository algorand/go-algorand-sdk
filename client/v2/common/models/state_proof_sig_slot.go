package models

// StateProofSigSlot defines a model for StateProofSigSlot.
type StateProofSigSlot struct {
	// LowerSigWeight (l) The total weight of signatures in the lower-numbered slots.
	LowerSigWeight uint64 `json:"lower-sig-weight,omitempty"`

	// Signature
	Signature StateProofSignature `json:"signature,omitempty"`
}
