package types

// MessageHash represents the message that a state proof will attest to.
type MessageHash [32]byte

// StateProofType identifies a particular configuration of state proofs.
type StateProofType uint64

// Message represents the message that the state proofs are attesting to. This message can be
// used by lightweight client and gives it the ability to verify proofs on the Algorand's state.
// In addition to that proof, this message also contains fields that
// are needed in order to verify the next state proofs (VotersCommitment and LnProvenWeight).
type Message struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
	// BlockHeadersCommitment contains a commitment on all light block headers within a state proof interval.
	BlockHeadersCommitment []byte `codec:"b,allocbound=Sha256Size"`
	VotersCommitment       []byte `codec:"v,allocbound=SumhashDigestSize"`
	LnProvenWeight         uint64 `codec:"P"`
	FirstAttestedRound     uint64 `codec:"f"`
	LastAttestedRound      uint64 `codec:"l"`
}

// StateProofTxnFields captures the fields used for stateproof transactions.
type StateProofTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	StateProofType StateProofType `codec:"sptype"`
	StateProof     interface{}    `codec:"sp"`
	Message        Message        `codec:"spmsg"`
}
