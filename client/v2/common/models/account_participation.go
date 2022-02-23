package models

// AccountParticipation accountParticipation describes the parameters used by this
// account in consensus protocol.
type AccountParticipation struct {
	// SelectionParticipationKey (sel) Selection public key (if any) currently
	// registered for this round.
	SelectionParticipationKey []byte `json:"selection-participation-key"`

	// StateProofKey (stprf) Root of the state proof key (if any)
	StateProofKey []byte `json:"state-proof-key,omitempty"`

	// VoteFirstValid (voteFst) First round for which this participation is valid.
	VoteFirstValid uint64 `json:"vote-first-valid"`

	// VoteKeyDilution (voteKD) Number of subkeys in each batch of participation keys.
	VoteKeyDilution uint64 `json:"vote-key-dilution"`

	// VoteLastValid (voteLst) Last round for which this participation is valid.
	VoteLastValid uint64 `json:"vote-last-valid"`

	// VoteParticipationKey (vote) root participation public key (if any) currently
	// registered for this round.
	VoteParticipationKey []byte `json:"vote-participation-key"`
}
