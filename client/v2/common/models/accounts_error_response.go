package models

// AccountsErrorResponse an error response for the /v2/accounts API endpoint, with
// optional information about limits that were exceeded.
type AccountsErrorResponse struct {
	// Address
	Address string `json:"address,omitempty"`

	// MaxResults
	MaxResults uint64 `json:"max-results,omitempty"`

	// Message
	Message string `json:"message"`

	// TotalAppsOptedIn
	TotalAppsOptedIn uint64 `json:"total-apps-opted-in,omitempty"`

	// TotalAssetsOptedIn
	TotalAssetsOptedIn uint64 `json:"total-assets-opted-in,omitempty"`

	// TotalCreatedApps
	TotalCreatedApps uint64 `json:"total-created-apps,omitempty"`

	// TotalCreatedAssets
	TotalCreatedAssets uint64 `json:"total-created-assets,omitempty"`
}
