package models

// Account account information at a given round.
// Definition:
// data/basics/userBalance.go : AccountData
type Account struct {
	// Address the account public key
	Address string `json:"address"`

	// Amount (algo) total number of MicroAlgos in the account
	Amount uint64 `json:"amount"`

	// AmountWithoutPendingRewards specifies the amount of MicroAlgos in the account,
	// without the pending rewards.
	AmountWithoutPendingRewards uint64 `json:"amount-without-pending-rewards"`

	// AppsLocalState (appl) applications local data stored in this account.
	// Note the raw object uses `map[int] -> AppLocalState` for this type.
	AppsLocalState []ApplicationLocalState `json:"apps-local-state,omitempty"`

	// AppsTotalExtraPages (teap) the sum of all extra application program pages for
	// this account.
	AppsTotalExtraPages uint64 `json:"apps-total-extra-pages,omitempty"`

	// AppsTotalSchema (tsch) stores the sum of all of the local schemas and global
	// schemas in this account.
	// Note: the raw account uses `StateSchema` for this type.
	AppsTotalSchema ApplicationStateSchema `json:"apps-total-schema,omitempty"`

	// Assets (asset) assets held by this account.
	// Note the raw object uses `map[int] -> AssetHolding` for this type.
	Assets []AssetHolding `json:"assets,omitempty"`

	// AuthAddr (spend) the address against which signing should be checked. If empty,
	// the address of the current account is used. This field can be updated in any
	// transaction by setting the RekeyTo field.
	AuthAddr string `json:"auth-addr,omitempty"`

	// ClosedAtRound round during which this account was most recently closed.
	ClosedAtRound uint64 `json:"closed-at-round,omitempty"`

	// CreatedApps (appp) parameters of applications created by this account including
	// app global data.
	// Note: the raw account uses `map[int] -> AppParams` for this type.
	CreatedApps []Application `json:"created-apps,omitempty"`

	// CreatedAssets (apar) parameters of assets created by this account.
	// Note: the raw account uses `map[int] -> Asset` for this type.
	CreatedAssets []Asset `json:"created-assets,omitempty"`

	// CreatedAtRound round during which this account first appeared in a transaction.
	CreatedAtRound uint64 `json:"created-at-round,omitempty"`

	// Deleted whether or not this account is currently closed.
	Deleted bool `json:"deleted,omitempty"`

	// Participation accountParticipation describes the parameters used by this account
	// in consensus protocol.
	Participation AccountParticipation `json:"participation,omitempty"`

	// PendingRewards amount of MicroAlgos of pending rewards in this account.
	PendingRewards uint64 `json:"pending-rewards"`

	// RewardBase (ebase) used as part of the rewards computation. Only applicable to
	// accounts which are participating.
	RewardBase uint64 `json:"reward-base,omitempty"`

	// Rewards (ern) total rewards of MicroAlgos the account has received, including
	// pending rewards.
	Rewards uint64 `json:"rewards"`

	// Round the round for which this information is relevant.
	Round uint64 `json:"round"`

	// SigType indicates what type of signature is used by this account, must be one
	// of:
	// * sig
	// * msig
	// * lsig
	// * or null if unknown
	SigType string `json:"sig-type,omitempty"`

	// Status (onl) delegation status of the account's MicroAlgos
	// * Offline - indicates that the associated account is delegated.
	// * Online - indicates that the associated account used as part of the delegation
	// pool.
	// * NotParticipating - indicates that the associated account is neither a
	// delegator nor a delegate.
	Status string `json:"status"`

	// TotalAppsOptedIn the count of all applications that have been opted in,
	// equivalent to the count of application local data (AppLocalState objects) stored
	// in this account.
	TotalAppsOptedIn uint64 `json:"total-apps-opted-in"`

	// TotalAssetsOptedIn the count of all assets that have been opted in, equivalent
	// to the count of AssetHolding objects held by this account.
	TotalAssetsOptedIn uint64 `json:"total-assets-opted-in"`

	// TotalBoxBytes for app-accounts only. The total number of bytes allocated for the
	// keys and values of boxes which belong to the associated application.
	TotalBoxBytes uint64 `json:"total-box-bytes"`

	// TotalBoxes for app-accounts only. The total number of boxes which belong to the
	// associated application.
	TotalBoxes uint64 `json:"total-boxes"`

	// TotalCreatedApps the count of all apps (AppParams objects) created by this
	// account.
	TotalCreatedApps uint64 `json:"total-created-apps"`

	// TotalCreatedAssets the count of all assets (AssetParams objects) created by this
	// account.
	TotalCreatedAssets uint64 `json:"total-created-assets"`
}
