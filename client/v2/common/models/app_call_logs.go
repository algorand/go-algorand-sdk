package models

// AppCallLogs the logged messages from an app call along with the app ID and outer
// transaction ID. Logs appear in the same order that they were emitted.
type AppCallLogs struct {
	// ApplicationIndex the application from which the logs were generated
	ApplicationIndex uint64 `json:"application-index"`

	// Logs an array of logs
	Logs [][]byte `json:"logs"`

	// Txid the transaction ID of the outer app call that lead to these logs
	Txid string `json:"txId"`
}
