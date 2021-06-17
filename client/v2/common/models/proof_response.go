package models

// ProofResponse proof of transaction in a block.
type ProofResponse struct {
	// Idx index of the transaction in the block's payset.
	Idx uint64 `json:"idx"`

	// Proof merkle proof of transaction membership.
	Proof []byte `json:"proof"`

	// Stibhash hash of SignedTxnInBlock for verifying proof.
	Stibhash []byte `json:"stibhash"`
}
