package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

const algodAuthHeader = "X-Algo-API-Token"

type Client common.Client

// get performs a GET request to the specific path against the server, assumes JSON response
func (c *Client) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Get(ctx, response, path, request, headers)
}

// getMsgpack performs a GET request to the specific path against the server, assumes msgpack response
func (c *Client) getMsgpack(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return (*common.Client)(c).GetRawMsgpack(ctx, response, path, request, headers)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (c *Client) post(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Post(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing a ClientV2 for a given endpoint.
func MakeClient(address string, apiToken string) (c *Client, err error) {
	commonClient, err := common.MakeClient(address, algodAuthHeader, apiToken)
	c = (*Client)(commonClient)
	return
}

func (c *Client) AccountInformation(account string) *AccountInformation {
	return &AccountInformation{c: c, account: account}
}

func (c *Client) Block(round uint64) *Block {
	return &Block{c: c, round: round}
}

func (c *Client) HealthCheck() *HealthCheck {
	return &HealthCheck{c: c}
}

func (c *Client) PendingTransactionInformation(txid string) *PendingTransactionInformation {
	return &PendingTransactionInformation{c: c, txid: txid}
}

func (c *Client) PendingTransactionsByAddress(address string) *PendingTransactionInformationByAddress {
	return &PendingTransactionInformationByAddress{c: c, address: address}
}

func (c *Client) PendingTransactions() *PendingTransactions {
	return &PendingTransactions{c: c}
}

func (c *Client) RegisterParticipationKeys(account string) *RegisterParticipationKeys {
	return &RegisterParticipationKeys{c: c, account: account}
}

func (c *Client) SendRawTransaction(tx []byte) *SendRawTransaction {
	return &SendRawTransaction{c: c, stx: tx}
}

func (c *Client) Shutdown() *Shutdown {
	return &Shutdown{c: c}
}

func (c *Client) StatusAfterBlock(round uint64) *StatusAfterBlock {
	return &StatusAfterBlock{c: c}
}

func (c *Client) Status() *Status {
	return &Status{c: c}
}

func (c *Client) SuggestedParams() *SuggestedParams {
	return &SuggestedParams{c: c}
}

func (c *Client) Supply() *Supply {
	return &Supply{c: c}
}

func (c *Client) Versions() *Versions {
	return &Versions{c: c}
}
