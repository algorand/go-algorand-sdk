package models

// ProofResponse proof of transaction in a block.
type ProofResponse struct {
	// Hashtype the type of hash function used to create the proof, must be one of:
	// * sumhash
	// * sha512_256
	Hashtype string `json:"hashtype,omitempty"`

	// Idx index of the transaction in the block's payset.
	Idx uint64 `json:"idx"`

	// Proof merkle proof of transaction membership.
	Proof []byte `json:"proof"`

	// Stibhash hash of SignedTxnInBlock for verifying proof.
	Stibhash []byte `json:"stibhash"`

	// Treedepth represents the depth of the tree that is being proven, i.e. the number
	// of edges from a leaf to the root.
	Treedepth uint64 `json:"treedepth"`
}
