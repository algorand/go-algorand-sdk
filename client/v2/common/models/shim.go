package models

// Types in this file are used to help convert between hand written code
// that doesn't quite match the spec, and code generated from the spec.

// NodeStatus is the algod status report.
type NodeStatus NodeStatusResponse

// PendingTransactionInfoResponse is the single pending transaction response.
type PendingTransactionInfoResponse PendingTransactionResponse

// HealthCheckResponse defines model for HealthCheckResponse.
type HealthCheckResponse HealthCheck

// Supply
type Supply SupplyResponse

// VersionBuild
type VersionBuild BuildVersion
