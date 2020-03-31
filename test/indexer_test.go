package main

import (
	"github.com/cucumber/godog"
)

func IndexerContext(s *godog.Suite) {
	s.Step(`^we make any LookupAssetBalances call$`, weMakeAnyLookupAssetBalancesCall)
	s.Step(`^the parsed LookupAssetBalances response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have address "([^"]*)" amount (\d+) and frozen state (\d+)$`, theParsedLookupAssetBalancesResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveAddressAmountAndFrozenState)
	s.Step(`^we make any LookupAssetTransactions call$`, weMakeAnyLookupAssetTransactionsCall)
	s.Step(`^the parsed LookupAssetTransactions response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedLookupAssetTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender)
	s.Step(`^we make any LookupAccountTransactions call$`, weMakeAnyLookupAccountTransactionsCall)
	s.Step(`^the parsed LookupAccountTransactions response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedLookupAccountTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender)
	s.Step(`^we make any LookupBlock call$`, weMakeAnyLookupBlockCall)
	s.Step(`^the parsed LookupBlock response should have proposer "([^"]*)"$`, theParsedLookupBlockResponseShouldHaveProposer)
	s.Step(`^we make any LookupAccountByID call$`, weMakeAnyLookupAccountByIDCall)
	s.Step(`^the parsed LookupAccountByID response should have address "([^"]*)"$`, theParsedLookupAccountByIDResponseShouldHaveAddress)
	s.Step(`^we make any LookupAssetByID call$`, weMakeAnyLookupAssetByIDCall)
	s.Step(`^the parsed LookupAssetByID response should have index (\d+)$`, theParsedLookupAssetByIDResponseShouldHaveIndex)
	s.Step(`^we make any SearchAccounts call$`, weMakeAnySearchAccountsCall)
	s.Step(`^the parsed SearchAccounts response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have address "([^"]*)"$`, theParsedSearchAccountsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAddress)
	s.Step(`^we make any SearchForTransactions call$`, weMakeAnySearchForTransactionsCall)
	s.Step(`^the parsed SearchForTransactions response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have sender "([^"]*)"$`, theParsedSearchForTransactionsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveSender)
	s.Step(`^we make any SearchForAssets call$`, weMakeAnySearchForAssetsCall)
	s.Step(`^the parsed SearchForAssets response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have asset index (\d+)$`, theParsedSearchForAssetsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAssetIndex)
	s.BeforeScenario(func(interface{}) {})
}

//
//func buildMockIndexerAndClient(jsonfile string) (*httptest.Server, indexer.Client, error) {
//	jsonBytes, err := loadMockJsons("", jsonfile)
//	if err != nil {
//		return nil, indexer.Client{}, err
//	}
//	mockIndexer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		w.Write(jsonBytes[0])
//	}))
//	noToken := ""
//	indexerClient, err := indexer.MakeClient(mockIndexer.URL, noToken)
//	return mockIndexer, indexerClient, err
//}
func mockHttpResponsesInLoadedFrom(arg1, arg2 string) error {
	return godog.ErrPending
}

func weMakeAnyLookupAssetBalancesCall() error {
	return godog.ErrPending
}

func theParsedLookupAssetBalancesResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveAddressAmountAndFrozenState(arg1, arg2, arg3 int, arg4 string, arg5, arg6 int) error {
	return godog.ErrPending
}

func weMakeAnyLookupAssetTransactionsCall() error {
	return godog.ErrPending
}

func theParsedLookupAssetTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender(arg1, arg2, arg3 int, arg4 string) error {
	return godog.ErrPending
}

func weMakeAnyLookupAccountTransactionsCall() error {
	return godog.ErrPending
}

func theParsedLookupAccountTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender(arg1, arg2, arg3 int, arg4 string) error {
	return godog.ErrPending
}

func weMakeAnyLookupBlockCall() error {
	return godog.ErrPending
}

func theParsedLookupBlockResponseShouldHaveProposer(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnyLookupAccountByIDCall() error {
	return godog.ErrPending
}

func theParsedLookupAccountByIDResponseShouldHaveAddress(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnyLookupAssetByIDCall() error {
	return godog.ErrPending
}

func theParsedLookupAssetByIDResponseShouldHaveIndex(arg1 int) error {
	return godog.ErrPending
}

func weMakeAnySearchAccountsCall() error {
	return godog.ErrPending
}

func theParsedSearchAccountsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAddress(arg1, arg2, arg3 int, arg4 string) error {
	return godog.ErrPending
}

func weMakeAnySearchForTransactionsCall() error {
	return godog.ErrPending
}

func theParsedSearchForTransactionsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveSender(arg1, arg2, arg3 int, arg4 string) error {
	return godog.ErrPending
}

func weMakeAnySearchForAssetsCall() error {
	return godog.ErrPending
}

func theParsedSearchForAssetsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAssetIndex(arg1, arg2, arg3, arg4 int) error {
	return godog.ErrPending
}
