package merklesignature

// MsgIsZero returns whether this is a zero value
func (z *Signature) MsgIsZero() bool {
	return ((*z).Signature.MsgIsZero()) && ((*z).VectorCommitmentIndex == 0) && ((*z).Proof.MsgIsZero()) && ((*z).VerifyingKey.MsgIsZero())
}
