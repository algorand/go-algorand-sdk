package models

// MerkleArrayProof defines a model for MerkleArrayProof.
type MerkleArrayProof struct {
	// HashFactory
	HashFactory HashFactory `json:"hash-factory,omitempty"`

	// Path (pth)
	Path [][]byte `json:"path,omitempty"`

	// TreeDepth (td)
	TreeDepth uint64 `json:"tree-depth,omitempty"`
}
