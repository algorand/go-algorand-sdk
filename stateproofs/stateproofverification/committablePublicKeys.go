package stateproofverification

import (
	"encoding/binary"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification/merklesignature"
)

type (
	// committablePublicKeyArray used to arrange the keys so a merkle tree could be build on them.
	//msgp:ignore committablePublicKeyArray
	committablePublicKeyArray struct {
		keys        []FalconSigner
		firstValid  uint64
		keyLifetime uint64
	}

	// CommittablePublicKey  is used to create a binary representation of public keys in the merkle
	// signature scheme.
	CommittablePublicKey struct {
		VerifyingKey FalconVerifier
		Round        uint64
	}
)

// ToBeHashed returns the sequence of bytes that would be used as an input for the hash function when creating a merkle tree.
// In order to create a more SNARK-friendly commitment we must avoid using the msgpack infrastructure.
// msgpack creates a compressed representation of the struct which might be varied in length, this will
// be bad for creating SNARK
func (e *CommittablePublicKey) ToBeHashed() (stateprooftypes.HashID, []byte) {
	verifyingRawKey := e.VerifyingKey.GetFixedLengthHashableRepresentation()

	roundAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(roundAsBytes, e.Round)

	schemeAsBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(schemeAsBytes, merklesignature.CryptoPrimitivesID)

	keyCommitment := make([]byte, 0, len(schemeAsBytes)+len(verifyingRawKey)+len(roundAsBytes))
	keyCommitment = append(keyCommitment, schemeAsBytes...)
	keyCommitment = append(keyCommitment, roundAsBytes...)
	keyCommitment = append(keyCommitment, verifyingRawKey...)

	return stateprooftypes.KeysInMSS, keyCommitment
}
