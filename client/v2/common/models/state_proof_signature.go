package models

// StateProofSignature defines a model for StateProofSignature.
type StateProofSignature struct {
	// FalconSignature
	FalconSignature []byte `json:"falcon-signature,omitempty"`

	// MerkleArrayIndex
	MerkleArrayIndex uint64 `json:"merkle-array-index,omitempty"`

	// Proof
	Proof MerkleArrayProof `json:"proof,omitempty"`

	// VerifyingKey (vkey)
	VerifyingKey []byte `json:"verifying-key,omitempty"`
}
