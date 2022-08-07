package stateproofverification

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
)

type committableSignatureSlot struct {
	sigCommit           sigslotCommit
	serializedSignature []byte
	isEmptySlot         bool
}

// ErrIndexOutOfBound returned when an index is out of the array's bound
var ErrIndexOutOfBound = errors.New("index is out of bound")

// committableSignatureSlotArray is used to create a merkle tree on the stateproof's
// signature array. it serializes the MSS signatures using a specific format
// state proof signature array.
//msgp:ignore committableSignatureSlotArray
type committableSignatureSlotArray []sigslot

func (sc committableSignatureSlotArray) Length() uint64 {
	return uint64(len(sc))
}

func (sc committableSignatureSlotArray) Marshal(pos uint64) (stateprooftypes.Hashable, error) {
	if pos >= uint64(len(sc)) {
		return nil, fmt.Errorf("%w: pos %d past end %d", ErrIndexOutOfBound, pos, len(sc))
	}

	signatureSlot, err := buildCommittableSignature(sc[pos].sigslotCommit)
	if err != nil {
		return nil, err
	}

	return signatureSlot, nil
}

func buildCommittableSignature(sigCommit sigslotCommit) (*committableSignatureSlot, error) {
	if sigCommit.Sig.Signature == nil {
		return &committableSignatureSlot{isEmptySlot: true}, nil
	}
	sigBytes, err := sigCommit.Sig.GetFixedLengthHashableRepresentation()
	if err != nil {
		return nil, err
	}
	return &committableSignatureSlot{sigCommit: sigCommit, serializedSignature: sigBytes, isEmptySlot: false}, nil
}

// ToBeHashed returns the sequence of bytes that would be used as an input for the hash function when creating a merkle tree.
// In order to create a more SNARK-friendly commitment we must avoid using the msgpack infrastructure.
// msgpack creates a compressed representation of the struct which might be varied in length, this will
// be bad for creating SNARK
func (cs *committableSignatureSlot) ToBeHashed() (stateprooftypes.HashID, []byte) {
	if cs.isEmptySlot {
		return stateprooftypes.StateProofSig, []byte{}
	}
	binaryLValue := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryLValue, cs.sigCommit.L)

	sigSlotByteRepresentation := make([]byte, 0, len(binaryLValue)+len(cs.serializedSignature))
	sigSlotByteRepresentation = append(sigSlotByteRepresentation, binaryLValue...)
	sigSlotByteRepresentation = append(sigSlotByteRepresentation, cs.serializedSignature...)

	return stateprooftypes.StateProofSig, sigSlotByteRepresentation
}
