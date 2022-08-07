package stateprooftypes

import "bytes"

//msgp:allocbound GenericDigest MaxHashDigestSize
type GenericDigest []byte

// ToSlice is used inside the Tree itself when interacting with TreeDigest
func (d GenericDigest) ToSlice() []byte { return d }

// IsEqual compare two digests
func (d GenericDigest) IsEqual(other GenericDigest) bool {
	return bytes.Equal(d, other)
}

// IsEmpty checks wether the generic digest is an empty one or not
func (d GenericDigest) IsEmpty() bool {
	return len(d) == 0
}
