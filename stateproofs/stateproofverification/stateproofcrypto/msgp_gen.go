package stateproofcrypto

// MsgIsZero returns whether this is a zero value
func (z *FalconVerifier) MsgIsZero() bool {
	return ((*z).PublicKey == (FalconPublicKey{}))
}

// MsgIsZero returns whether this is a zero value
func (z FalconSignature) MsgIsZero() bool {
	return len(z) == 0
}
