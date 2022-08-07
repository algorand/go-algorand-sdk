package stateproofverification

import (
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
)

// A Layer of the Merkle tree consists of a dense array of hashes at that
// level of the tree.  Hashes beyond the end of the array (e.g., if the
// number of leaves is not an exact power of 2) are implicitly zero.
//msgp:allocbound Layer MaxNumLeavesOnEncodedTree
type Layer []stateprooftypes.GenericDigest

// A pair represents an internal node in the Merkle tree.
type pair struct {
	l              stateprooftypes.GenericDigest
	r              stateprooftypes.GenericDigest
	hashDigestSize int
}

func (p *pair) ToBeHashed() (stateprooftypes.HashID, []byte) {
	// hashing of internal node will always be fixed length.
	// If one of the children is missing we use [0...0].
	// The size of the slice is based on the relevant hash function output size
	buf := make([]byte, 2*p.hashDigestSize)
	copy(buf[:], p.l[:])
	copy(buf[len(p.l):], p.r[:])
	return stateprooftypes.MerkleArrayNode, buf[:]
}
