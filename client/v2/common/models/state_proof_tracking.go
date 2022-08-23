package models

// StateProofTracking defines a model for StateProofTracking.
type StateProofTracking struct {
	// NextRound (n) Next round for which we will accept a state proof transaction.
	NextRound uint64 `json:"next-round,omitempty"`

	// OnlineTotalWeight (t) The total number of microalgos held by the online accounts
	// during the StateProof round.
	OnlineTotalWeight uint64 `json:"online-total-weight,omitempty"`

	// Type state Proof Type. Note the raw object uses map with this as key.
	Type uint64 `json:"type,omitempty"`

	// VotersCommitment (v) Root of a vector commitment containing online accounts that
	// will help sign the proof.
	VotersCommitment []byte `json:"voters-commitment,omitempty"`
}
