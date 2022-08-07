package stateproofverification

import (
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
)

// Proof is used to convince a verifier about membership of leaves: h0,h1...hn
// at indexes i0,i1...in on a tree. The verifier has a trusted value of the tree
// root hash.
type Proof struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Path is bounded by MaxNumLeavesOnEncodedTree since there could be multiple reveals, and
	// given the distribution of the elt positions and the depth of the tree,
	// the path length can increase up to 2^MaxEncodedTreeDepth / 2
	// TODO: allocbound
	Path        []stateprooftypes.GenericDigest `codec:"pth,allocbound=MaxNumLeavesOnEncodedTree/2"`
	HashFactory stateprooftypes.HashFactory     `codec:"hsh"`
	// TreeDepth represents the depth of the tree that is being proven.
	// It is the number of edges from the root to a leaf.
	TreeDepth uint8 `codec:"td"`
}

// SingleLeafProof is used to convince a verifier about membership of a specific
// leaf h at index i on a tree. The verifier has a trusted value of the tree
// root hash. it corresponds to merkle verification path.
type SingleLeafProof struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Proof
}

// GetFixedLengthHashableRepresentation serializes the proof into a sequence of bytes.
// it basically concatenates the elements of the verification path one after another.
// The function returns a fixed length array for each hash function. which is 1 + MaxEncodedTreeDepth * digestsize
//
// the path is guaranteed to be less than MaxEncodedTreeDepth and if the path length is less
// than MaxEncodedTreeDepth, array will have leading zeros (to fill the array to MaxEncodedTreeDepth * digestsize).
// more details could be found in the Algorand's spec.
func (p *SingleLeafProof) GetFixedLengthHashableRepresentation() []byte {
	hash := p.HashFactory.NewHash()

	var binProof = make([]byte, 0, 1+(MaxEncodedTreeDepth*hash.Size()))

	proofLenByte := p.TreeDepth
	binProof = append(binProof, proofLenByte)

	zeroDigest := make([]byte, hash.Size())

	for i := uint8(0); i < (MaxEncodedTreeDepth - proofLenByte); i++ {
		binProof = append(binProof, zeroDigest...)
	}

	for i := uint8(0); i < proofLenByte; i++ {
		if i < proofLenByte && p.Path[i] != nil {
			binProof = append(binProof, p.Path[i]...)
		} else {
			binProof = append(binProof, zeroDigest...)
		}
	}

	return binProof
}

// ToProof export a Proof from a SingleProof. The result is
// used as an input for merklearray.Verify or merklearray.VerifyVectorCommitment
func (p *SingleLeafProof) ToProof() *Proof {
	return &p.Proof
}
