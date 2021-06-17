package models

// BlockRewards fields relating to rewards,
type BlockRewards struct {
	// FeeSink (fees) accepts transaction fees, it can only spend to the incentive
	// pool.
	FeeSink string `json:"fee-sink"`

	// RewardsCalculationRound (rwcalr) number of leftover MicroAlgos after the
	// distribution of rewards-rate MicroAlgos for every reward unit in the next round.
	RewardsCalculationRound uint64 `json:"rewards-calculation-round"`

	// RewardsLevel (earn) How many rewards, in MicroAlgos, have been distributed to
	// each RewardUnit of MicroAlgos since genesis.
	RewardsLevel uint64 `json:"rewards-level"`

	// RewardsPool (rwd) accepts periodic injections from the fee-sink and continually
	// redistributes them as rewards.
	RewardsPool string `json:"rewards-pool"`

	// RewardsRate (rate) Number of new MicroAlgos added to the participation stake
	// from rewards at the next round.
	RewardsRate uint64 `json:"rewards-rate"`

	// RewardsResidue (frac) Number of leftover MicroAlgos after the distribution of
	// RewardsRate/rewardUnits MicroAlgos for every reward unit in the next round.
	RewardsResidue uint64 `json:"rewards-residue"`
}
