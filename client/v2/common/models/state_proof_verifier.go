package models

// StateProofVerifier defines a model for StateProofVerifier.
type StateProofVerifier struct {
	// Commitment (cmt) Represents the root of the vector commitment tree.
	Commitment []byte `json:"commitment,omitempty"`

	// KeyLifetime (lf) Key lifetime.
	KeyLifetime uint64 `json:"key-lifetime,omitempty"`
}
