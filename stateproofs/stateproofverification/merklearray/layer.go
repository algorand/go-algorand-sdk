package merklearray

import (
	"github.com/algorand/go-algorand-sdk/stateproofs/transactionverificationtypes"
)

// A Layer of the Merkle tree consists of a dense array of hashes at that
// level of the tree.  Hashes beyond the end of the array (e.g., if the
// number of leaves is not an exact power of 2) are implicitly zero.
//msgp:allocbound Layer MaxNumLeavesOnEncodedTree
type Layer []transactionverificationtypes.GenericDigest

// A pair represents an internal node in the Merkle tree.
type pair struct {
	l              transactionverificationtypes.GenericDigest
	r              transactionverificationtypes.GenericDigest
	hashDigestSize int
}

func (p *pair) ToBeHashed() (transactionverificationtypes.HashID, []byte) {
	// hashing of internal node will always be fixed length.
	// If one of the children is missing we use [0...0].
	// The size of the slice is based on the relevant hash function output size
	buf := make([]byte, 2*p.hashDigestSize)
	copy(buf[:], p.l[:])
	copy(buf[len(p.l):], p.r[:])
	return transactionverificationtypes.MerkleArrayNode, buf[:]
}
