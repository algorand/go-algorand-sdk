package algod

import (
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

// /v2/applications/{application-id}
// Given a application id, it returns application information including creator,
// approval and clear programs, global and local schemas, and global state.
func (c *Client) GetApplicationByID(applicationId uint64) *GetApplicationByID {
	return &GetApplicationByID{c: c, applicationId: applicationId}
}

// /v2/assets/{asset-id}
// Given a asset id, it returns asset information including creator, name, total
// supply and special addresses.
func (c *Client) GetAssetByID(assetId uint64) *GetAssetByID {
	return &GetAssetByID{c: c, assetId: assetId}
}

// /v2/teal/compile
// Given TEAL source code in plain text, return base64 encoded program bytes and
// base32 SHA512_256 hash of program bytes (Address style).
func (c *Client) TealCompile(source []byte) *TealCompile {
	return &TealCompile{c: c, source: source}
}

// /v2/teal/dryrun
// Executes TEAL program(s) in context and returns debugging information about the
// execution.
func (c *Client) TealDryrun(request models.DryrunRequest) *TealDryrun {
	return &TealDryrun{c: c, request: request}
}
