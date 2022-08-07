package stateproofverification

import (
	"fmt"
	"math/bits"
)

// merkleTreeToVectorCommitmentIndex Translate an index of an element on a merkle tree to an index on the vector commitment.
// The given index must be within the range of the elements in the tree (assume this number is 1^pathLen)
func merkleTreeToVectorCommitmentIndex(msbIndex uint64, pathLen uint8) (uint64, error) {
	if msbIndex >= (1 << pathLen) {
		return 0, fmt.Errorf("msbIndex %d >= 1^pathLen %d: %w", msbIndex, 1<<pathLen, ErrPosOutOfBound)
	}
	return bits.Reverse64(msbIndex) >> (64 - pathLen), nil
}
