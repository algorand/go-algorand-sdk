package models

// StateProofFields (sp) represents a state proof.
// Definition:
// crypto/stateproof/structs.go : StateProof
type StateProofFields struct {
	// PartProofs (P)
	PartProofs MerkleArrayProof `json:"part-proofs,omitempty"`

	// PositionsToReveal (pr) Sequence of reveal positions.
	PositionsToReveal []uint64 `json:"positions-to-reveal,omitempty"`

	// Reveals (r) Note that this is actually stored as a map[uint64] - Reveal in the
	// actual msgp
	Reveals []StateProofReveal `json:"reveals,omitempty"`

	// SaltVersion (v) Salt version of the merkle signature.
	SaltVersion uint64 `json:"salt-version,omitempty"`

	// SigCommit (c)
	SigCommit []byte `json:"sig-commit,omitempty"`

	// SigProofs (S)
	SigProofs MerkleArrayProof `json:"sig-proofs,omitempty"`

	// SignedWeight (w)
	SignedWeight uint64 `json:"signed-weight,omitempty"`
}
