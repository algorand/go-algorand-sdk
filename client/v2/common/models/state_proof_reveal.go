package models

// StateProofReveal defines a model for StateProofReveal.
type StateProofReveal struct {
	// Participant (p)
	Participant StateProofParticipant `json:"participant,omitempty"`

	// Position the position in the signature and participants arrays corresponding to
	// this entry.
	Position uint64 `json:"position,omitempty"`

	// SigSlot (s)
	SigSlot StateProofSigSlot `json:"sig-slot,omitempty"`
}
