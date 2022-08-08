package stateprooftypes

// MsgIsZero returns whether this is a zero value
func (z *HashFactory) MsgIsZero() bool {
	return ((*z).HashType == 0)
}
