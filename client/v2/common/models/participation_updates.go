package models

// ParticipationUpdates participation account data that needs to be checked/acted
// on by the network.
type ParticipationUpdates struct {
	// ExpiredParticipationAccounts (partupdrmv) a list of online accounts that needs
	// to be converted to offline since their participation key expired.
	ExpiredParticipationAccounts []string `json:"expired-participation-accounts,omitempty"`
}
