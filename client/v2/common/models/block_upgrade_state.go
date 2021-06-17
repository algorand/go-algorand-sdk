package models

// BlockUpgradeState fields relating to a protocol upgrade.
type BlockUpgradeState struct {
	// CurrentProtocol (proto) The current protocol version.
	CurrentProtocol string `json:"current-protocol"`

	// NextProtocol (nextproto) The next proposed protocol version.
	NextProtocol string `json:"next-protocol,omitempty"`

	// NextProtocolApprovals (nextyes) Number of blocks which approved the protocol
	// upgrade.
	NextProtocolApprovals uint64 `json:"next-protocol-approvals,omitempty"`

	// NextProtocolSwitchOn (nextswitch) Round on which the protocol upgrade will take
	// effect.
	NextProtocolSwitchOn uint64 `json:"next-protocol-switch-on,omitempty"`

	// NextProtocolVoteBefore (nextbefore) Deadline round for this protocol upgrade (No
	// votes will be consider after this round).
	NextProtocolVoteBefore uint64 `json:"next-protocol-vote-before,omitempty"`
}
