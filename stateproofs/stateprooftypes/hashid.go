package stateprooftypes

// HashID is a domain separation prefix for an object type that might be hashed
// This ensures, for example, the hash of a transaction will never collide with the hash of a vote
type HashID string

// Hash IDs for specific object types, in lexicographic order.
// Hash IDs must be PREFIX-FREE (no hash ID is a prefix of another).
const (
	BlockHeader256  HashID = "B256"
	MerkleArrayNode HashID = "MA"

	StateProofMessage HashID = "spm"

	TxnMerkleLeaf HashID = "TL"
)
