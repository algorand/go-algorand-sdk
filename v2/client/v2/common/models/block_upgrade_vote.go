package models

// BlockUpgradeVote fields relating to voting for a protocol upgrade.
type BlockUpgradeVote struct {
	// UpgradeApprove (upgradeyes) Indicates a yes vote for the current proposal.
	UpgradeApprove bool `json:"upgrade-approve,omitempty"`

	// UpgradeDelay (upgradedelay) Indicates the time between acceptance and execution.
	UpgradeDelay uint64 `json:"upgrade-delay,omitempty"`

	// UpgradePropose (upgradeprop) Indicates a proposed upgrade.
	UpgradePropose string `json:"upgrade-propose,omitempty"`
}
