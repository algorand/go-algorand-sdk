package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

// TODO ejr received peer feedback to have Block return types.Block not models.Block?

func (client Client) Shutdown(ctx context.Context, timeout models.ShutdownParams, headers ...*common.Header) error {
	return client.post(ctx, nil, "/shutdown", timeout, headers)
}

func (client Client) RegisterParticipationKeys(ctx context.Context, account string, params models.RegisterParticipationKeysAccountIdParams, headers ...*common.Header) error {
	return client.post(ctx, nil, fmt.Sprintf("/register-participation-keys/%s", account), nil, headers)
}

func (client Client) PendingTransactionInformation(ctx context.Context, txid string, params models.GetPendingTransactionsParams, headers ...*common.Header) (result types.Transaction, err error) {
	if params.Format == "json" {
		var response models.Transaction
		err = client.get(ctx, &response, fmt.Sprintf("/transactions/pending/%s", txid), params, headers)
		// TODO built result from response
	} else if params.Format == "msgpack" {
		err = client.get(ctx, &result, fmt.Sprintf("/transactions/pending/%s", txid), params, headers)
	} else {
		err = fmt.Errorf("unrecognized format %s, valid formats are json or msgpack", params.Format)
	}
	return
}

func (client Client) SendRawTransaction(ctx context.Context, txBytes []byte, headers ...*common.Header) (txid string, err error) {
	var response models.TxId
	headers = append(headers, &common.Header{Key: "Content-Type", Value: "application/x-binary"})
	err = client.post(ctx, &response, "/transactions", nil, headers)
	txid = string(response)
	return
}

func (client Client) PendingTransactionsByAddress(ctx context.Context, account string, params models.GetPendingTransactionsByAddressParams, headers ...*common.Header) (result []types.Transaction, err error) {
	if params.Format == "json" {
		var response []models.Transaction
		err = client.get(ctx, &response, fmt.Sprintf("/accounts/%s/transactions/pending", account), params, headers)
		// TODO built result from response
	} else if params.Format == "msgpack" {
		err = client.get(ctx, &result, fmt.Sprintf("/accounts/%s/transactions/pending", account), params, headers)
	} else {
		err = fmt.Errorf("unrecognized format %s, valid formats are json or msgpack", params.Format)
	}
	return
}

func (client Client) Status(ctx context.Context, headers ...*common.Header) (status models.NodeStatus, err error) {
	err = client.get(ctx, &status, "/status", nil, headers)
	return
}

func (client Client) Supply(ctx context.Context, headers ...*common.Header) (supply models.Supply, err error) {
	err = client.get(ctx, &supply, "/ledger/supply", nil, headers)
	return
}

func (client Client) StatusAfterBlock(ctx context.Context, round uint64, headers ...*common.Header) (status models.NodeStatus, err error) {
	err = client.get(ctx, &status, fmt.Sprintf("/status/wait-for-block-after/%d", round), nil, headers)
	return
}

func (client Client) AccountInformation(ctx context.Context, address string, headers ...*common.Header) (result models.Account, err error) {
	err = client.get(ctx, &result, fmt.Sprintf("/accounts/%s", address), nil, headers)
	return
}

func (client Client) Block(ctx context.Context, round uint64, params models.GetBlockParams, headers ...*common.Header) (result models.Block, err error) {
	err = client.get(ctx, &result, fmt.Sprintf("/block/%d", round), params, headers)
	if params.Format == "json" {
		var response models.RawBlockJson
		err = client.get(ctx, &response, fmt.Sprintf("/block/%d", round), params, headers)
		// TODO built result from response
	} else if params.Format == "msgpack" {
		var response models.RawBlockMsgpack
		err = client.get(ctx, &response, fmt.Sprintf("/block/%d", round), params, headers)
		// TODO built result from response
	} else {
		err = fmt.Errorf("unrecognized format %s, valid formats are json or msgpack", params.Format)
	}
	return
}

func (client Client) SuggestedParams(ctx context.Context, headers ...*common.Header) (params types.SuggestedParams, err error) {
	var response models.TransactionParams
	err = client.get(ctx, &response, "/transactions/params", nil, headers)
	params = types.SuggestedParams{
		Fee:              types.MicroAlgos(response.Fee),
		GenesisID:        response.GenesisID,
		GenesisHash:      response.Genesishash,
		FirstRoundValid:  types.Round(response.LastRound),
		LastRoundValid:   types.Round(response.LastRound + 1000),
		ConsensusVersion: response.ConsensusVersion,
	}
	return
}

// Versions retrieves the VersionResponse from the running node
// the VersionResponse includes data like version number and genesis ID
func (client Client) Versions(ctx context.Context, headers ...*common.Header) (response models.Version, err error) {
	err = client.get(ctx, &response, "/versions", nil, headers)
	return
}

// HealthCheck does a health check on the the potentially running node,
// returning an error if the API is down
func (client Client) HealthCheck(ctx context.Context, headers ...*common.Header) error {
	return client.get(ctx, nil, "/health", nil, headers)
}
