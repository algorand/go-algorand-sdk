package models

// AccountErrorResponse an error response for the AccountInformation endpoint, with
// optional information about limits that were exceeded.
type AccountErrorResponse struct {
	// Data
	Data string `json:"data,omitempty"`

	// MaxResults
	MaxResults uint64 `json:"max-results,omitempty"`

	// Message
	Message string `json:"message"`

	// TotalAppsLocalState
	TotalAppsLocalState uint64 `json:"total-apps-local-state,omitempty"`

	// TotalAssets
	TotalAssets uint64 `json:"total-assets,omitempty"`

	// TotalCreatedApps
	TotalCreatedApps uint64 `json:"total-created-apps,omitempty"`

	// TotalCreatedAssets
	TotalCreatedAssets uint64 `json:"total-created-assets,omitempty"`
}
