package models

// IndexerStateProofMessage defines a model for IndexerStateProofMessage.
type IndexerStateProofMessage struct {
	// BlockHeadersCommitment (b)
	BlockHeadersCommitment []byte `json:"block-headers-commitment,omitempty"`

	// FirstAttestedRound (f)
	FirstAttestedRound uint64 `json:"first-attested-round,omitempty"`

	// LatestAttestedRound (l)
	LatestAttestedRound uint64 `json:"latest-attested-round,omitempty"`

	// LnProvenWeight (P)
	LnProvenWeight uint64 `json:"ln-proven-weight,omitempty"`

	// VotersCommitment (v)
	VotersCommitment []byte `json:"voters-commitment,omitempty"`
}
