package types

// Hash IDs for specific object types, in lexicographic order.
// Hash IDs must be PREFIX-FREE (no hash ID is a prefix of another).
const (
	BlockHeader256 HashID = "B256"

	StateProofMessage HashID = "spm"

	TxnMerkleLeaf HashID = "TL"
)
