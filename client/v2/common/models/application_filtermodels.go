package models

// SearchForApplicationsParams defines parameters for SearchForApplications
type SearchForApplicationsParams struct {
	// ApplicationId application ID
	ApplicationId uint64 `url:"application-id,omitempty"`

	// Limit maximum number of results to return.
	Limit uint64 `url:"limit,omitempty"`

	// Next the next page of results. Use the next token provided by the previous
	// results.
	Next string `url:"next,omitempty"`
}
