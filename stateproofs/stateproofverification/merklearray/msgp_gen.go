package merklearray

// MsgIsZero returns whether this is a zero value
func (z *SingleLeafProof) MsgIsZero() bool {
	return (len((*z).Proof.Path) == 0) && ((*z).Proof.HashFactory.MsgIsZero()) && ((*z).Proof.TreeDepth == 0)
}
