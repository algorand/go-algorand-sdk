package indexer

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func (client Client) LookupAssetBalances(ctx context.Context, assetIndex uint64, params models.LookupAssetBalancesParams,
	headers ...*common.Header) (validRound uint64, holders []models.MiniAssetHolding, err error) {
	var response models.AssetBalancesResponse
	err = client.get(ctx, &response, fmt.Sprintf("/assets/%d/balances", assetIndex), params, headers)
	validRound = response.CurrentRound
	holders = response.Balances
	return
}

func (client Client) LookupAssetTransactions(ctx context.Context, assetIndex uint64, params models.LookupAssetTransactionsParams,
	headers ...*common.Header) (validRound uint64, transactions []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = client.get(ctx, &response, fmt.Sprintf("/assets/%d/transactions", assetIndex), params, headers)
	validRound = response.CurrentRound
	transactions = response.Transactions
	return
}

func (client Client) LookupAccountTransactions(ctx context.Context, account string, params models.LookupAccountTransactionsParams,
	headers ...*common.Header) (validRound uint64, transactions []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = client.get(ctx, &response, fmt.Sprintf("/accounts/%s/transactions", account), params, headers)
	validRound = response.CurrentRound
	transactions = response.Transactions
	return
}

func (client Client) LookupBlock(ctx context.Context, round uint64, headers ...*common.Header) (block models.Block, err error) {
	var response models.BlockResponse
	err = client.get(ctx, &response, fmt.Sprintf("/blocks/%d", round), nil, headers)
	block = models.Block(response)
	return
}

func (client Client) LookupAccountByID(ctx context.Context, account string, params models.LookupAccountByIDParams, headers ...*common.Header) (validRound uint64, result models.Account, err error) {
	var response struct {
		CurrentRound uint64         `json:"current-round"`
		Account      models.Account `json:"account"`
	}
	err = client.get(ctx, &response, fmt.Sprintf("/accounts/%s", account), params, headers)
	validRound = response.CurrentRound
	result = response.Account
	return
}

func (client Client) LookupAssetByID(ctx context.Context, assetIndex uint64, headers ...*common.Header) (validRound uint64, result models.Asset, err error) {
	var response struct {
		CurrentRound uint64       `json:"current-round"`
		Asset        models.Asset `json:"asset"`
	}
	err = client.get(ctx, &response, fmt.Sprintf("/assets/%d", assetIndex), nil, headers)
	validRound = response.CurrentRound
	result = response.Asset
	return
}

func (client Client) SearchAccounts(ctx context.Context, params models.SearchAccountsParams, headers ...*common.Header) (validRound uint64, result []models.Account, err error) {
	var response models.AccountsResponse
	err = client.get(ctx, &response, "/accounts", params, headers)
	validRound = response.CurrentRound
	result = response.Accounts
	return
}

func (client Client) SearchForTransactions(ctx context.Context, params models.SearchForTransactionsParams, headers ...*common.Header) (validRound uint64, result []models.Transaction, err error) {
	var response models.TransactionsResponse
	err = client.get(ctx, &response, "/transactions", params, headers)
	validRound = response.CurrentRound
	result = response.Transactions
	return
}

func (client Client) SearchForAssets(ctx context.Context, params models.SearchForAssetsParams, headers ...*common.Header) (validRound uint64, result []models.Asset, err error) {
	var response models.AssetsResponse
	err = client.get(ctx, &response, "/assets", params, headers)
	validRound = response.CurrentRound
	result = response.Assets
	return
}
