package models

// TransactionKeyreg fields for a keyreg transaction.
// Definition:
// data/transactions/keyreg.go : KeyregTxnFields
type TransactionKeyreg struct {
	// NonParticipation (nonpart) Mark the account as participating or
	// non-participating.
	NonParticipation bool `json:"non-participation,omitempty"`

	// SelectionParticipationKey (selkey) Public key used with the Verified Random
	// Function (VRF) result during committee selection.
	SelectionParticipationKey []byte `json:"selection-participation-key,omitempty"`

	// StateProofKey (sprfkey) State proof key used in key registration transactions.
	StateProofKey []byte `json:"state-proof-key,omitempty"`

	// VoteFirstValid (votefst) First round this participation key is valid.
	VoteFirstValid uint64 `json:"vote-first-valid,omitempty"`

	// VoteKeyDilution (votekd) Number of subkeys in each batch of participation keys.
	VoteKeyDilution uint64 `json:"vote-key-dilution,omitempty"`

	// VoteLastValid (votelst) Last round this participation key is valid.
	VoteLastValid uint64 `json:"vote-last-valid,omitempty"`

	// VoteParticipationKey (votekey) Participation public key used in key registration
	// transactions.
	VoteParticipationKey []byte `json:"vote-participation-key,omitempty"`
}
