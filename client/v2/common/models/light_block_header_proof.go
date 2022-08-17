package models

// LightBlockHeaderProof proof of membership and position of a light block header.
type LightBlockHeaderProof struct {
	// Index the index of the light block header in the vector commitment tree
	Index uint64 `json:"index"`

	// Proof the encoded proof.
	Proof []byte `json:"proof"`

	// Treedepth represents the depth of the tree that is being proven, i.e. the number
	// of edges from a leaf to the root.
	Treedepth uint64 `json:"treedepth"`
}
