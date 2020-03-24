package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/algorand/go-algorand-sdk/v2client/common/models"
	"github.com/algorand/go-algorand-sdk/v2client/indexer"

	"github.com/cucumber/godog"
)

func IndexerContext(s *godog.Suite) {
	s.Step(`^we make any LookupAssetBalances call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnyLookupAssetBalancesCallReturnMockResponse)
	s.Step(`^the parsed LookupAssetBalances response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have address "([^"]*)" amount (\d+) and frozen state (\d+)$`, theParsedLookupAssetBalancesResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveAddressAmountAndFrozenState)
	s.Step(`^we make any LookupAssetTransactions call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnyLookupAssetTransactionsCallReturnMockResponse)
	s.Step(`^the parsed LookupAssetTransactions response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedLookupAssetTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender)
	s.Step(`^we make any LookupAccountTransactions call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnyLookupAccountTransactionsCallReturnMockResponse)
	s.Step(`^the parsed LookupAccountTransactions response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedLookupAccountTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender)
	s.Step(`^we make any LookupBlock call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnyLookupBlockCallReturnMockResponse)
	s.Step(`^the parsed LookupBlock response should have round "([^"]*)"$`, theParsedLookupBlockResponseShouldHaveRound)
	s.Step(`^we make any LookupAccountByID call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnyLookupAccountByIDCallReturnMockResponse)
	s.Step(`^the parsed LookupAccountByID response should have address "([^"]*)"$`, theParsedLookupAccountByIDResponseShouldHaveAddress)
	s.Step(`^we make any LookupAssetByID call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnyLookupAssetByIDCallReturnMockResponse)
	s.Step(`^the parsed LookupAssetByID response should have index (\d+)$`, theParsedLookupAssetByIDResponseShouldHaveIndex)
	s.Step(`^we make any SearchAccounts call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnySearchAccountsCallReturnMockResponse)
	s.Step(`^the parsed SearchAccounts response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have address "([^"]*)"$`, theParsedSearchAccountsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAddress)
	s.Step(`^we make any SearchForTransactions call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnySearchForTransactionsCallReturnMockResponse)
	s.Step(`^the parsed SearchForTransactions response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have sender "([^"]*)"$`, theParsedSearchForTransactionsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveSender)
	s.Step(`^we make any SearchForAssets call, return mock response "([^"]*)" and expect error string to contain "([^"]*)$`, weMakeAnySearchForAssetsCallReturnMockResponse)
	s.Step(`^the parsed SearchForAssets response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have asset index (\d+)$`, theParsedSearchForAssetsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAssetIndex)

	s.BeforeScenario(func(interface{}) {
	})
}

func buildMockIndexerAndClient(jsonfile string) (*httptest.Server, indexer.Client, error) {
	jsonBytes, err := loadMockJson(jsonfile)
	if err != nil {
		return nil, indexer.Client{}, err
	}
	mockIndexer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}))
	noToken := ""
	indexerClient, err := indexer.MakeClient(mockIndexer.URL, noToken)
	return mockIndexer, indexerClient, err
}

var validRound uint64
var assetHolders []models.MiniAssetHolding

func weMakeAnyLookupAssetBalancesCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, assetHolders, err = indexerClient.LookupAssetBalances(context.Background(), 0, models.LookupAssetBalancesParams{}, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedLookupAssetBalancesResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveAddressAmountAndFrozenState(roundNum, length, idx int, address string, amount, frozenState int) error {
	if uint64(roundNum) != validRound {
		return fmt.Errorf("roundNum %d mismatched with validRound %d", roundNum, validRound)
	}
	if len(assetHolders) != length {
		return fmt.Errorf("len(assetHolders) %d mismatched with length %d", len(assetHolders), length)
	}
	examinedHolder := assetHolders[idx]
	var expectedFrozenState bool
	if frozenState == 0 {
		expectedFrozenState = false
	} else {
		expectedFrozenState = true
	}
	if examinedHolder.IsFrozen != expectedFrozenState {
		return fmt.Errorf("examinedHolder.IsFrozen %v mismatched with expectedFrozenState %v", examinedHolder.IsFrozen, expectedFrozenState)
	}
	if examinedHolder.Address != address {
		return fmt.Errorf("examinedHolder.Address %s mismatched with expected address %s", examinedHolder.Address, address)
	}
	if examinedHolder.Amount != uint64(amount) {
		return fmt.Errorf("examinedHolder.Amount %d mismatched with expected amount %d", examinedHolder.Amount, uint64(amount))
	}
	return nil
}

var transactions []models.Transaction

func weMakeAnyLookupAssetTransactionsCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, transactions, err = indexerClient.LookupAssetTransactions(context.Background(), 0, models.LookupAssetTransactionsParams{}, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedLookupAssetTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender(roundNum, length, idx int, sender string) error {
	if uint64(roundNum) != validRound {
		return fmt.Errorf("roundNum %d mismatched with validRound %d", roundNum, validRound)
	}
	if len(transactions) != length {
		return fmt.Errorf("len(transactions) %d mismatched with length %d", len(transactions), length)
	}
	examinedTransaction := transactions[idx]
	if examinedTransaction.Sender != sender {
		return fmt.Errorf("examinedTransaction.Sender %s mismatched with expected sender %s", examinedTransaction.Sender, sender)
	}
	return nil
}

func weMakeAnyLookupAccountTransactionsCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, transactions, err = indexerClient.LookupAccountTransactions(context.Background(), "", models.LookupAccountTransactionsParams{}, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedLookupAccountTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender(roundNum, length, idx int, sender string) error {
	if uint64(roundNum) != validRound {
		return fmt.Errorf("roundNum %d mismatched with validRound %d", roundNum, validRound)
	}
	if len(transactions) != length {
		return fmt.Errorf("len(transactions) %d mismatched with length %d", len(transactions), length)
	}
	examinedTransaction := transactions[idx]
	if examinedTransaction.Sender != sender {
		return fmt.Errorf("examinedTransaction.Sender %s mismatched with expected sender %s", examinedTransaction.Sender, sender)
	}
	return nil
}

var foundBlock models.Block

func weMakeAnyLookupBlockCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	foundBlock, err = indexerClient.LookupBlock(context.Background(), 0, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedLookupBlockResponseShouldHaveRound(round int) error {
	if foundBlock.Round != uint64(round) {
		return fmt.Errorf("block round %d mismatch with expected round %d", foundBlock.Round, round)
	}
	return nil
}

var foundAccount models.Account

func weMakeAnyLookupAccountByIDCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, foundAccount, err = indexerClient.LookupAccountByID(context.Background(), "", models.LookupAccountByIDParams{}, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedLookupAccountByIDResponseShouldHaveAddress(address string) error {
	if foundAccount.Address != address {
		return fmt.Errorf("lookup account address %s mismatch with expected address %s", foundAccount.Address, address)
	}
	return nil
}

var foundAsset models.Asset

func weMakeAnyLookupAssetByIDCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, foundAsset, err = indexerClient.LookupAssetByID(context.Background(), 0, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedLookupAssetByIDResponseShouldHaveIndex(idx int) error {
	if foundAsset.Index != uint64(idx) {
		return fmt.Errorf("lookup asset index %d mismatch with expected index %d", foundAsset.Index, uint64(idx))
	}
	return nil
}

var foundAccounts []models.Account

func weMakeAnySearchAccountsCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, foundAccounts, err = indexerClient.SearchAccounts(context.Background(), models.SearchAccountsParams{}, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedSearchAccountsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAddress(roundNum, length, idx int, address string) error {
	if uint64(roundNum) != validRound {
		return fmt.Errorf("roundNum %d mismatched with validRound %d", roundNum, validRound)
	}
	if len(foundAccounts) != length {
		return fmt.Errorf("len(foundAccounts) %d mismatched with length %d", len(foundAccounts), length)
	}
	examinedAccount := foundAccounts[idx]
	if examinedAccount.Address != address {
		return fmt.Errorf("examinedAccount.Address %s mismatched with expected address %s", examinedAccount.Address, address)
	}
	return nil
}

func weMakeAnySearchForTransactionsCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, transactions, err = indexerClient.SearchForTransactions(context.Background(), models.SearchForTransactionsParams{}, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedSearchForTransactionsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveSender(roundNum, length, idx int, sender string) error {
	if uint64(roundNum) != validRound {
		return fmt.Errorf("roundNum %d mismatched with validRound %d", roundNum, validRound)
	}
	if len(transactions) != length {
		return fmt.Errorf("len(transactions) %d mismatched with length %d", len(transactions), length)
	}
	examinedTransaction := transactions[idx]
	if examinedTransaction.Sender != sender {
		return fmt.Errorf("examinedTransaction.Sender %s mismatched with expected sender %s", examinedTransaction.Sender, sender)
	}
	return nil
}

var foundAssets []models.Asset

func weMakeAnySearchForAssetsCallReturnMockResponse(jsonfile string, errorString string) error {
	mockIndexer, indexerClient, err := buildMockIndexerAndClient(jsonfile)
	if mockIndexer != nil {
		defer mockIndexer.Close()
	}
	if err != nil {
		return err
	}
	validRound, foundAssets, err = indexerClient.SearchForAssets(context.Background(), models.SearchForAssetsParams{}, nil)
	return confirmErrorContainsString(err, errorString)
}

func theParsedSearchForAssetsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAssetIndex(roundNum, length, idx, assetIndex int) error {
	if uint64(roundNum) != validRound {
		return fmt.Errorf("roundNum %d mismatched with validRound %d", roundNum, validRound)
	}
	if len(foundAssets) != length {
		return fmt.Errorf("len(foundAssets) %d mismatched with length %d", len(foundAssets), length)
	}
	examinedAsset := foundAssets[idx]
	if examinedAsset.Index != uint64(assetIndex) {
		return fmt.Errorf("examinedAsset.Index %d mismatched with expected index %d", examinedAsset.Index, uint64(assetIndex))
	}
	return nil
}
