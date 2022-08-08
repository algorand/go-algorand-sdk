package stateproof

import (
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification/merklearray"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification/merklesignature"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification/stateproofbasics"
	"github.com/algorand/go-algorand-sdk/stateproofs/transactionverificationtypes"
)

//msgp:ignore sigslot
type sigslot struct {
	// Weight is the weight of the participant signing this message.
	// This information is tracked here for convenience, but it does
	// not appear in the commitment to the sigs array; it comes from
	// the Weight field of the corresponding participant.
	Weight uint64

	// Include the parts of the sigslot that form the commitment to
	// the sigs array.
	sigslotCommit
}

// A sigslotCommit is a single slot in the sigs array that forms the state proof.
type sigslotCommit struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Sig is a signature by the participant on the expected message.
	Sig merklesignature.Signature `codec:"s"`

	// L is the total weight of signatures in lower-numbered slots.
	// This is initialized once the builder has collected a sufficient
	// number of signatures.
	L uint64 `codec:"l"`
}

// Reveal is a single array position revealed as part of a state
// proof.  It reveals an element of the signature array and
// the corresponding element of the participants array.
type Reveal struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	SigSlot sigslotCommit                `codec:"s"`
	Part    stateproofbasics.Participant `codec:"p"`
}

// StateProof represents a proof on Algorand's state.
type StateProof struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	SigCommit                  transactionverificationtypes.GenericDigest `codec:"c"`
	SignedWeight               uint64                                     `codec:"w"`
	SigProofs                  merklearray.Proof                          `codec:"S"`
	PartProofs                 merklearray.Proof                          `codec:"P"`
	MerkleSignatureSaltVersion byte                                       `codec:"v"`
	// Reveals is a sparse map from the position being revealed
	// to the corresponding elements from the sigs and participants
	// arrays.
	Reveals           map[uint64]Reveal `codec:"r,allocbound=MaxReveals"`
	PositionsToReveal []uint64          `codec:"pr,allocbound=MaxReveals"`
}
