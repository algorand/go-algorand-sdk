package indexer

// /v2/applications/{application-id}
// Lookup application.
func (c *Client) LookupApplicationByID(applicationId uint64) *LookupApplicationByID {
	return &LookupApplicationByID{c: c, applicationId: applicationId}
}

// /v2/applications
// Search for applications
func (c *Client) SearchForApplications() *SearchForApplications {
	return &SearchForApplications{c: c}
}

