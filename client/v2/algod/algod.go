package algod

import (
	"context"
	"net/http"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
)

const authHeader = "X-Algo-API-Token"

type Client common.Client

// delete performs a DELETE request to the specific path against the server, assumes JSON response
func (c *Client) delete(ctx context.Context, response interface{}, path string, body interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Delete(ctx, response, path, body, headers)
}

// get performs a GET request to the specific path against the server, assumes JSON response
func (c *Client) get(ctx context.Context, response interface{}, path string, body interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Get(ctx, response, path, body, headers)
}

// getMsgpack performs a GET request to the specific path against the server, assumes msgpack response
func (c *Client) getMsgpack(ctx context.Context, response interface{}, path string, body interface{}, headers []*common.Header) error {
	return (*common.Client)(c).GetRawMsgpack(ctx, response, path, body, headers)
}

// getMsgpack performs a GET request to the specific path against the server, assumes msgpack response
func (c *Client) getRaw(ctx context.Context, path string, body interface{}, headers []*common.Header) ([]byte, error) {
	return (*common.Client)(c).GetRaw(ctx, path, body, headers)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (c *Client) post(ctx context.Context, response interface{}, path string, params interface{}, headers []*common.Header, body interface{}) error {
	return (*common.Client)(c).Post(ctx, response, path, params, headers, body)
}

// MakeClient is the factory for constructing a ClientV2 for a given endpoint.
func MakeClient(address string, apiToken string) (c *Client, err error) {
	commonClient, err := common.MakeClient(address, authHeader, apiToken)
	c = (*Client)(commonClient)
	return
}

// MakeClientWithHeaders is the factory for constructing a ClientV2 for a
// given endpoint with custom headers.
func MakeClientWithHeaders(address string, apiToken string, headers []*common.Header) (c *Client, err error) {
	commonClientWithHeaders, err := common.MakeClientWithHeaders(address, authHeader, apiToken, headers)
	c = (*Client)(commonClientWithHeaders)
	return
}

// MakeClientWithTransport is the factory for constructing a Client for a given endpoint with a
// custom HTTP Transport as well as optional additional user defined headers.
func MakeClientWithTransport(address string, apiToken string, headers []*common.Header, transport http.RoundTripper) (c *Client, err error) {
	commonClientWithTransport, err := common.MakeClientWithTransport(address, authHeader, apiToken, headers, transport)
	c = (*Client)(commonClientWithTransport)
	return
}

func (c *Client) HealthCheck() *HealthCheck {
	return &HealthCheck{c: c}
}

func (c *Client) GetReady() *GetReady {
	return &GetReady{c: c}
}

func (c *Client) GetGenesis() *GetGenesis {
	return &GetGenesis{c: c}
}

func (c *Client) Versions() *Versions {
	return &Versions{c: c}
}

func (c *Client) AccountInformation(address string) *AccountInformation {
	return &AccountInformation{c: c, address: address}
}

func (c *Client) AccountAssetInformation(address string, assetId uint64) *AccountAssetInformation {
	return &AccountAssetInformation{c: c, address: address, assetId: assetId}
}

func (c *Client) AccountApplicationInformation(address string, applicationId uint64) *AccountApplicationInformation {
	return &AccountApplicationInformation{c: c, address: address, applicationId: applicationId}
}

func (c *Client) PendingTransactionsByAddress(address string) *PendingTransactionsByAddress {
	return &PendingTransactionsByAddress{c: c, address: address}
}

func (c *Client) Block(round uint64) *Block {
	return &Block{c: c, round: round}
}

func (c *Client) GetBlockHash(round uint64) *GetBlockHash {
	return &GetBlockHash{c: c, round: round}
}

func (c *Client) GetTransactionProof(round uint64, txid string) *GetTransactionProof {
	return &GetTransactionProof{c: c, round: round, txid: txid}
}

func (c *Client) Supply() *Supply {
	return &Supply{c: c}
}

func (c *Client) Status() *Status {
	return &Status{c: c}
}

func (c *Client) StatusAfterBlock(round uint64) *StatusAfterBlock {
	return &StatusAfterBlock{c: c, round: round}
}

func (c *Client) SendRawTransaction(rawtxn []byte) *SendRawTransaction {
	return &SendRawTransaction{c: c, rawtxn: rawtxn}
}

func (c *Client) SimulateTransaction(request models.SimulateRequest) *SimulateTransaction {
	return &SimulateTransaction{c: c, request: request}
}

func (c *Client) SuggestedParams() *SuggestedParams {
	return &SuggestedParams{c: c}
}

func (c *Client) PendingTransactions() *PendingTransactions {
	return &PendingTransactions{c: c}
}

func (c *Client) PendingTransactionInformation(txid string) *PendingTransactionInformation {
	return &PendingTransactionInformation{c: c, txid: txid}
}

func (c *Client) GetLedgerStateDelta(round uint64) *GetLedgerStateDelta {
	return &GetLedgerStateDelta{c: c, round: round}
}

func (c *Client) GetTransactionGroupLedgerStateDeltasForRound(round uint64) *GetTransactionGroupLedgerStateDeltasForRound {
	return &GetTransactionGroupLedgerStateDeltasForRound{c: c, round: round}
}

func (c *Client) GetLedgerStateDeltaForTransactionGroup(id string) *GetLedgerStateDeltaForTransactionGroup {
	return &GetLedgerStateDeltaForTransactionGroup{c: c, id: id}
}

func (c *Client) GetStateProof(round uint64) *GetStateProof {
	return &GetStateProof{c: c, round: round}
}

func (c *Client) GetLightBlockHeaderProof(round uint64) *GetLightBlockHeaderProof {
	return &GetLightBlockHeaderProof{c: c, round: round}
}

func (c *Client) GetApplicationByID(applicationId uint64) *GetApplicationByID {
	return &GetApplicationByID{c: c, applicationId: applicationId}
}

func (c *Client) GetApplicationBoxes(applicationId uint64) *GetApplicationBoxes {
	return &GetApplicationBoxes{c: c, applicationId: applicationId}
}

func (c *Client) GetApplicationBoxByName(applicationId uint64, name []byte) *GetApplicationBoxByName {
	return (&GetApplicationBoxByName{c: c, applicationId: applicationId}).name(name)
}

func (c *Client) GetAssetByID(assetId uint64) *GetAssetByID {
	return &GetAssetByID{c: c, assetId: assetId}
}

func (c *Client) UnsetSyncRound() *UnsetSyncRound {
	return &UnsetSyncRound{c: c}
}

func (c *Client) GetSyncRound() *GetSyncRound {
	return &GetSyncRound{c: c}
}

func (c *Client) SetSyncRound(round uint64) *SetSyncRound {
	return &SetSyncRound{c: c, round: round}
}

func (c *Client) TealCompile(source []byte) *TealCompile {
	return &TealCompile{c: c, source: source}
}

func (c *Client) TealDisassemble(source []byte) *TealDisassemble {
	return &TealDisassemble{c: c, source: source}
}

func (c *Client) TealDryrun(request models.DryrunRequest) *TealDryrun {
	return &TealDryrun{c: c, request: request}
}

func (c *Client) GetBlockTimeStampOffset() *GetBlockTimeStampOffset {
	return &GetBlockTimeStampOffset{c: c}
}

func (c *Client) SetBlockTimeStampOffset(offset uint64) *SetBlockTimeStampOffset {
	return &SetBlockTimeStampOffset{c: c, offset: offset}
}

func (c *Client) BlockRaw(round uint64) *BlockRaw {
	return &BlockRaw{c: c, round: round}
}
