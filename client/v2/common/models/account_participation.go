package models

// AccountParticipation accountParticipation describes the parameters used by this
// account in consensus protocol.
type AccountParticipation struct {
	// SelectionParticipationKey selection public key (if any) currently registered for
	// this round.
	SelectionParticipationKey []byte `json:"selection-participation-key"`

	// StateProofKey root of the state proof key (if any)
	StateProofKey []byte `json:"state-proof-key,omitempty"`

	// VoteFirstValid first round for which this participation is valid.
	VoteFirstValid uint64 `json:"vote-first-valid"`

	// VoteKeyDilution number of subkeys in each batch of participation keys.
	VoteKeyDilution uint64 `json:"vote-key-dilution"`

	// VoteLastValid last round for which this participation is valid.
	VoteLastValid uint64 `json:"vote-last-valid"`

	// VoteParticipationKey root participation public key (if any) currently registered
	// for this round.
	VoteParticipationKey []byte `json:"vote-participation-key"`
}
