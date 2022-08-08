package transactionverificationtypes

import (
	"crypto/sha256"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
)

// EncodedStateProof represents the msgpack encoded state proof.
type EncodedStateProof []byte

// MessageHash represents the message that a state proof will attest to.
type MessageHash [32]byte

// Message represents the message that the state proofs are attesting to. This message can be
// used by lightweight client and gives it the ability to verify proofs on the Algorand's state.
// In addition to that proof, this message also contains fields that
// are needed in order to verify the next state proofs (VotersCommitment and LnProvenWeight).
type Message struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
	// BlockHeadersCommitment contains a commitment on all light block headers within a state proof interval.
	BlockHeadersCommitment []byte `codec:"b,allocbound=Sha256Size"`
	VotersCommitment       []byte `codec:"v,allocbound=MaxHashDigestSize"`
	LnProvenWeight         uint64 `codec:"P"`
	FirstAttestedRound     uint64 `codec:"f"`
	LastAttestedRound      uint64 `codec:"l"`
}

// ToBeHashed returns the bytes of the message.
func (m Message) ToBeHashed() (HashID, []byte) {
	return StateProofMessage, msgpack.Encode(&m)
}

// IntoStateProofMessageHash returns a hashed representation fitting the state proof messages.
func (m Message) IntoStateProofMessageHash() MessageHash {
	digest := GenericHashObj(sha256.New(), m)
	result := MessageHash{}
	copy(result[:], digest)
	return result
}
