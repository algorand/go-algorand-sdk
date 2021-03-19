package models

// Types in this file are used to help convert between hand written code
// that doesn't quite match the spec, and code generated from the spec.

// NodeStatus is the algod status report.
type NodeStatus NodeStatusResponse

// PendingTransactionInfoResponse is the single pending transaction response.
type PendingTransactionInfoResponse PendingTransactionResponse

type Supply SupplyResponse

//type VersionBuild BuildVersion

// VersionBuild defines model for the current algod build version information.
type VersionBuild struct {
	Branch      string `json:"branch"`
	BuildNumber uint64 `json:"build-number"`
	Channel     string `json:"channel"`
	CommitHash  []byte `json:"commit-hash"`
	Major       uint64 `json:"major"`
	Minor       uint64 `json:"minor"`
}
