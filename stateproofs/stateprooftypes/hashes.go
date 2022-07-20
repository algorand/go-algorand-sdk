package stateprooftypes

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

// HashType represents different hash functions
type HashType uint16

// TODO: Add error object
// Validate verifies that the hash type is in a valid range.
func (h HashType) Validate() error {
	if h >= MaxHashType {
		return fmt.Errorf("error invalid object")
	}
	return nil
}

// types of hashes
const (
	Sha512_256 HashType = iota
	Sumhash
	Sha256
	MaxHashType
)

type Hashable interface {
	ToBeHashed() (HashID, []byte)
}

// HashRep appends the correct hashid before the message to be hashed.
func HashRep(h Hashable) []byte {
	hashid, data := h.ToBeHashed()
	return append([]byte(hashid), data...)
}

// Sumhash512DigestSize  The size in bytes of the sumhash checksum
const Sumhash512DigestSize = 64

// MaxHashDigestSize is used to bound the max digest size. it is important to change it if a hash with
// a longer output is introduced.
const MaxHashDigestSize = SumhashDigestSize

//size of each hash
const (
	Sha512_256Size    = sha512.Size256
	SumhashDigestSize = Sumhash512DigestSize
	Sha256Size        = sha256.Size
)

// TODO: add error
func UnmarshalHashFunc(hashStr string) (hash.Hash, error) {
	switch hashStr {
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash function detected")
	}
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
