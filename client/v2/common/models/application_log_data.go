package models

// ApplicationLogData stores the global information associated with an application.
type ApplicationLogData struct {
	// Logs logs for the application being executed by the transaction.
	Logs [][]byte `json:"logs"`

	// Txid transaction ID
	Txid string `json:"txid"`
}
