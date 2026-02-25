package models

// AccountApplicationsInformationResponse accountApplicationsInformationResponse
// contains a list of application resources for an account.
type AccountApplicationsInformationResponse struct {
	// ApplicationResources
	ApplicationResources []AccountApplicationResource `json:"application-resources,omitempty"`

	// NextToken used for pagination, when making another request provide this token
	// with the next parameter. The next token is the next application ID to use as the
	// pagination cursor.
	NextToken string `json:"next-token,omitempty"`

	// Round the round for which this information is relevant.
	Round uint64 `json:"round"`
}
