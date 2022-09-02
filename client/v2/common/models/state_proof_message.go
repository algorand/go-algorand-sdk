package models

// StateProofMessage represents the message that the state proofs are attesting to.
type StateProofMessage struct {
	// Blockheaderscommitment the vector commitment root on all light block headers
	// within a state proof interval.
	Blockheaderscommitment []byte `json:"BlockHeadersCommitment"`

	// Firstattestedround the first round the message attests to.
	Firstattestedround uint64 `json:"FirstAttestedRound"`

	// Lastattestedround the last round the message attests to.
	Lastattestedround uint64 `json:"LastAttestedRound"`

	// Lnprovenweight an integer value representing the natural log of the proven
	// weight with 16 bits of precision. This value would be used to verify the next
	// state proof.
	Lnprovenweight uint64 `json:"LnProvenWeight"`

	// Voterscommitment the vector commitment root of the top N accounts to sign the
	// next StateProof.
	Voterscommitment []byte `json:"VotersCommitment"`
}
