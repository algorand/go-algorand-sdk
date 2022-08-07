package merklesignature

import "github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"

// HashType/ hashSize relate to the type of hash this package uses.
const (
	MerkleSignatureSchemeHashFunction = stateprooftypes.Sumhash
	MerkleSignatureSchemeRootSize     = stateprooftypes.SumhashDigestSize
	// KeyLifetimeDefault defines the default lifetime of a key in the merkle signature scheme (in rounds).
	KeyLifetimeDefault = 256

	// SchemeSaltVersion is the current salt version of merkleSignature
	SchemeSaltVersion = byte(0)

	// CryptoPrimitivesID is an identification that the Merkle Signature Scheme uses a subset sum hash function
	// and a falcon signature scheme.
	CryptoPrimitivesID = uint16(0)
)
