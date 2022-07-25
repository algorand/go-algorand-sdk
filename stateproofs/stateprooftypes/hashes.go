package stateprooftypes

import (
	"hash"
)

type Hashable interface {
	ToBeHashed() (HashID, []byte)
}

// HashRep appends the correct hashid before the message to be hashed.
func HashRep(h Hashable) []byte {
	hashid, data := h.ToBeHashed()
	return append([]byte(hashid), data...)
}

func HashBytes(hash hash.Hash, m []byte) []byte {
	hash.Reset()
	hash.Write(m)
	outhash := hash.Sum(nil)
	return outhash
}

func GenericHashObj(hsh hash.Hash, h Hashable) []byte {
	rep := HashRep(h)
	return HashBytes(hsh, rep)
}
