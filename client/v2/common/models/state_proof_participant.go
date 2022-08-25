package models

// StateProofParticipant defines a model for StateProofParticipant.
type StateProofParticipant struct {
	// Verifier (p)
	Verifier StateProofVerifier `json:"verifier,omitempty"`

	// Weight (w)
	Weight uint64 `json:"weight,omitempty"`
}
