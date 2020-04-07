package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

type Client common.Client

// get performs a GET request to the specific path against the server
func (c *Client) get(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Get(ctx, response, path, request, headers)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (c *Client) post(ctx context.Context, response interface{}, path string, request interface{}, headers []*common.Header) error {
	return (*common.Client)(c).Post(ctx, response, path, request, headers)
}

// MakeClient is the factory for constructing a ClientV2 for a given endpoint.
func MakeClient(address string, apiToken string) (c *Client, err error) {
	commonClient, err := common.MakeClient(address, apiToken)
	c = (*Client)(commonClient)
	return
}

func (c *Client) NewAccountInformationService(account string) *AccountInformationService {
	return &AccountInformationService{c: c, account: account}
}

func (c *Client) NewBlockService(round uint64) *BlockService {
	return &BlockService{c: c, round: round}
}

func (c *Client) NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{c: c}
}

func (c *Client) NewPendingTransactionInformationService(txid string) *PendingTransactionInformationService {
	return &PendingTransactionInformationService{c: c, txid: txid}
}

func (c *Client) NewPendingTransactionsByAddressService(address string) *PendingTransactionInformationByAddressService {
	return &PendingTransactionInformationByAddressService{c: c, address: address}
}

func (c *Client) NewPendingTransactionsService() *PendingTransactionsService {
	return &PendingTransactionsService{c: c}
}

func (c *Client) NewRegisterParticipationKeysService(account string) *RegisterParticipationKeysService {
	return &RegisterParticipationKeysService{c: c, account: account}
}

func (c *Client) NewSendRawTransactionService(tx []byte) *SendRawTransactionService {
	return &SendRawTransactionService{c: c, stx: tx}
}

func (c *Client) NewShutdownService() *ShutdownService {
	return &ShutdownService{c: c}
}

func (c *Client) NewStatusAfterBlockService(round uint64) *StatusAfterBlockService {
	return &StatusAfterBlockService{c: c}
}

func (c *Client) NewStatusService() *StatusService {
	return &StatusService{c: c}
}

func (c *Client) NewSuggestedParamsService() *SuggestedParamsService {
	return &SuggestedParamsService{c: c}
}

func (c *Client) NewSupplyService() *SupplyService {
	return &SupplyService{c: c}
}

func (c *Client) NewVersionsService() *VersionsService {
	return &VersionsService{c: c}
}
