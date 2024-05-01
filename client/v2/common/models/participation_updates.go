package models

// ParticipationUpdates participation account data that needs to be checked/acted
// on by the network.
type ParticipationUpdates struct {
	// AbsentParticipationAccounts (partupabs) a list of online accounts that need to
	// be suspended.
	AbsentParticipationAccounts []string `json:"absent-participation-accounts,omitempty"`

	// ExpiredParticipationAccounts (partupdrmv) a list of online accounts that needs
	// to be converted to offline since their participation key expired.
	ExpiredParticipationAccounts []string `json:"expired-participation-accounts,omitempty"`
}
