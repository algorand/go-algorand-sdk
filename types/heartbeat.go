package types

// HeartbeatTxnFields captures the fields used for an account to prove it is
// online (really, it proves that an entity with the account's part keys is able
// to submit transactions, so it should be able to propose/vote.)
type HeartbeatTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// HbAddress is the account this txn is proving onlineness for.
	HbAddress Address `codec:"a"`

	// HbProof is a signature using HeartbeatAddress's partkey, thereby showing it is online.
	HbProof HeartbeatProof `codec:"prf"`

	// The final three fields are included to allow early, concurrent check of
	// the HbProof.

	// HbSeed must be the block seed for this transaction's firstValid block. It
	// is the message that must be signed with HbAddress's part key.
	HbSeed Seed `codec:"sd"`

	// HbVoteID must match the HbAddress account's current VoteID.
	HbVoteID VotePK `codec:"vid"`

	// HbKeyDilution must match HbAddress account's current KeyDilution.
	HbKeyDilution uint64 `codec:"kd"`

	// HbChallengeDiscount requests the challenge fee discount: when set, the
	// required fee is reduced by one min fee. It is optional even for a
	// challenged account (an account willing to pay the normal fee can leave it
	// off), so it is a request, not an assertion. apply verifies HbAddress is
	// actually under challenge before granting it. The flag is only allowed
	// once transaction size pricing is enabled (proto.TxnSizePricingEnabled());
	// it makes sense to think in terms of transaction fields changing fees
	// now, so it needs no separate consensus flag. Before then, the discount
	// was inferred from an underpaid singleton heartbeat instead (see
	// wellFormed and apply).
	HbChallengeDiscount bool `codec:"c"`
}
