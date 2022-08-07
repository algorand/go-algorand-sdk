package stateproofverification

import (
	"encoding/binary"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
)

// A Participant corresponds to an account whose AccountData.Status
// is Online, and for which the expected sigRound satisfies
// AccountData.VoteFirstValid <= sigRound <= AccountData.VoteLastValid.
//
// In the Algorand ledger, it is possible for multiple accounts to have
// the same PK.  Thus, the PK is not necessarily unique among Participants.
// However, each account will produce a unique Participant struct, to avoid
// potential DoS attacks where one account claims to have the same VoteID PK
// as another account.
type Participant struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// PK is the identifier used to verify the signature for a specific participant
	PK Verifier `codec:"p"`

	// Weight is AccountData.MicroAlgos.
	Weight uint64 `codec:"w"`
}

// ToBeHashed implements the crypto.Hashable interface.
// In order to create a more SNARK-friendly commitments on the signature we must avoid using the msgpack infrastructure.
// msgpack creates a compressed representation of the struct which might be varied in length, which will
// be bad for creating SNARK
func (p Participant) ToBeHashed() (stateprooftypes.HashID, []byte) {

	weightAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(weightAsBytes, p.Weight)

	keyLifetimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyLifetimeBytes, p.PK.KeyLifetime)

	publicKeyBytes := p.PK.Commitment

	partCommitment := make([]byte, 0, len(weightAsBytes)+len(publicKeyBytes)+len(keyLifetimeBytes))
	partCommitment = append(partCommitment, weightAsBytes...)
	partCommitment = append(partCommitment, keyLifetimeBytes...)
	partCommitment = append(partCommitment, publicKeyBytes[:]...)

	return stateprooftypes.StateProofPart, partCommitment
}
