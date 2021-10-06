package models

// ApplicationLogsResponse
type ApplicationLogsResponse struct {
	// ApplicationId (appidx) application index.
	ApplicationId uint64 `json:"application-id"`

	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// LogData
	LogData []ApplicationLogData `json:"log-data,omitempty"`

	// NextToken used for pagination, when making another request provide this token
	// with the next parameter.
	NextToken string `json:"next-token,omitempty"`
}
