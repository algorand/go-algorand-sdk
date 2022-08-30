package types

import "hash"

// GenericDigest is a digest that implements CustomSizeDigest, and can be used as hash output.
//msgp:allocbound GenericDigest MaxHashDigestSize
type GenericDigest []byte

// HashID is a domain separation prefix for an object type that might be hashed
// This ensures, for example, the hash of a transaction will never collide with the hash of a vote
type HashID string

// Hashable is an interface implemented by an object that can be represented
// with a sequence of bytes to be hashed or signed, together with a type ID
// to distinguish different types of objects.
type Hashable interface {
	ToBeHashed() (HashID, []byte)
}

// HashRep appends the correct hashid before the message to be hashed.
func HashRep(h Hashable) []byte {
	hashid, data := h.ToBeHashed()
	return append([]byte(hashid), data...)
}

// GenericHashObj Makes it easier to sum using hash interface and Hashable interface
func GenericHashObj(hsh hash.Hash, h Hashable) []byte {
	rep := HashRep(h)
	return hashBytes(hsh, rep)
}

func hashBytes(hash hash.Hash, m []byte) []byte {
	hash.Reset()
	hash.Write(m)
	outhash := hash.Sum(nil)
	return outhash
}
