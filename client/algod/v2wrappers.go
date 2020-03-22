package algod

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/types"
)

// TODO ejr handling of token

func (client Client) Shutdown(ctx context.Context, timeout models.GetV2ShutdownParams, headers ...*Header) error {
	// TODO EJR how to handle private security token
	return client.post(ctx, nil, "/shutdown", timeout, headers)
}

func (client Client) RegisterParticipationKeys(ctx context.Context, account string, params models.GetV2RegisterParticipationKeysAccountIdParams, headers ...*Header) error {
	// TODO EJR how to handle private security token
	return client.post(ctx, nil, fmt.Sprintf("/register-participation-keys/%s", account), nil, headers)
}

func (client Client) PendingTransactionInformation(ctx context.Context, txid string, params models.GetPendingTransactionsParams, headers ...*Header) (result types.Transaction, err error) {
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

func (client Client) SendRawTransaction(ctx context.Context, txBytes []byte, headers ...*Header) (txid string, err error) {
	var response models.TxId
	err = client.post(ctx, &response, "/transactions", nil, headers)
	txid = string(response)
	return
}

func (client Client) PendingTransactionsByAddress(ctx context.Context, account string, params models.GetPendingTransactionsByAddressParams, headers ...*Header) (result []types.Transaction, err error) {
	err = client.get(ctx, &result, fmt.Sprintf("/accounts/%s/transactions/pending", account), params, headers)
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

func (client Client) Status(ctx context.Context, headers ...*Header) (status models.NodeStatus, err error) {
	err = client.get(ctx, &status, "/status", nil, headers)
	return
}

func (client Client) Supply(ctx context.Context, headers ...*Header) (supply models.Supply, err error) {
	err = client.get(ctx, &supply, "/ledger/supply", nil, headers)
	return
}

func (client Client) StatusAfterBlock(ctx context.Context, round uint64, headers ...*Header) (status models.NodeStatus, err error) {
	err = client.get(ctx, &status, fmt.Sprintf("/statuswait-for-block-after/%d", round), nil, headers)
	return
}

func (client Client) AccountInformation(ctx context.Context, address string, headers ...*Header) (result models.Account, err error) {
	err = client.get(ctx, &result, fmt.Sprintf("/accounts/%s", address), nil, headers)
	return
}

func (client Client) Block(ctx context.Context, round uint64, params models.GetBlockParams, headers ...*Header) (result models.Block, err error) {
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

func (client Client) SuggestedParams(ctx context.Context, headers ...*Header) (params types.SuggestedParams, err error) {
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
