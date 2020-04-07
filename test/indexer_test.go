package main

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/cucumber/godog"
	"time"
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
	s.Step(`^we make a Lookup Asset Balances call against asset index (\d+) with limit <limit> afterAddress "([^"]*)" round (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+)$`, weMakeALookupAssetBalancesCallAgainstAssetIndexWithLimitLimitAfterAddressRoundCurrencyGreaterThanCurrencyLessThan)
	s.Step(`^we make a Lookup Asset Transactions call against asset index (\d+) with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime (\d+) afterTime (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) address "([^"]*)" addressRole "([^"]*)" ExcluseCloseTo "([^"]*)"$`, weMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseTo)
	s.Step(`^we make a Lookup Account Transactions call against account "([^"]*)" with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime (\d+) afterTime (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) addressRole "([^"]*)" ExcluseCloseTo "([^"]*)"$`, weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseTo)
	s.Step(`^we make a Lookup Block call against round (\d+)$`, weMakeALookupBlockCallAgainstRound)
	s.Step(`^we make a Lookup Account by ID call against account "([^"]*)" with round (\d+)$`, weMakeALookupAccountByIDCallAgainstAccountWithRound)
	s.Step(`^we make a Lookup Asset by ID call against asset index (\d+)$`, weMakeALookupAssetByIDCallAgainstAssetIndex)
	s.Step(`^we make a Search Accounts call with assetID (\d+) limit (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) and afterAddress "([^"]*)"$`, weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanAndAfterAddress)
	s.Step(`^we make a Search For Transactions call with account "([^"]*)" NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime (\d+) afterTime (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) addressRole "([^"]*)" ExcluseCloseTo "([^"]*)"$`, weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseTo)
	s.Step(`^we make a SearchForAssets call with limit (\d+) creator "([^"]*)" name "([^"]*)" unit "([^"]*)" index (\d+) and afterAsset (\d+)$`, weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndexAndAfterAsset)
	s.BeforeScenario(func(interface{}) {
		globalErrForExamination = nil
	})
}

var responseValidRound uint64
var assetBalancesResponse []models.MiniAssetHolding

func weMakeAnyLookupAssetBalancesCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, assetBalancesResponse, globalErrForExamination = indexerClient.NewLookupAssetBalancesService(0).Do(context.Background())
	return nil
}

func theParsedLookupAssetBalancesResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveAddressAmountAndFrozenState(round, length, idx int, address string, amount int, frozenState string) error {
	if responseValidRound != uint64(round) {
		return fmt.Errorf("response round %d did not match expected round %d", responseValidRound, uint64(round))
	}
	realLen := len(assetBalancesResponse)
	if realLen != length {
		return fmt.Errorf("response length %d did not match expected length %d", realLen, length)
	}
	scrutinizedElement := assetBalancesResponse[idx]
	if scrutinizedElement.Address != address {
		return fmt.Errorf("response address %s did not match expected address %s", scrutinizedElement.Address, address)
	}
	if scrutinizedElement.Amount != uint64(amount) {
		return fmt.Errorf("response amount %d did not match expected amount %d", scrutinizedElement.Amount, amount)
	}
	var isFrozenBool bool
	if frozenState == "false" {
		isFrozenBool = false
	} else if frozenState == "true" {
		isFrozenBool = true
	} else {
		return fmt.Errorf("unrecognized expected isFrozen value %s, allowed values are \"true\" \"false\"", frozenState)
	}
	if scrutinizedElement.IsFrozen != isFrozenBool {
		return fmt.Errorf("response frozen state %v did not match expected frozen state %v", scrutinizedElement.IsFrozen, isFrozenBool)
	}
	return nil
}

var assetTransactionsResponse []models.Transaction

func weMakeAnyLookupAssetTransactionsCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, assetTransactionsResponse, globalErrForExamination = indexerClient.NewLookupAssetTransactionsService(0).Do(context.Background())
	return nil
}

func theParsedLookupAssetTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender(round, length, idx int, sender string) error {
	if responseValidRound != uint64(round) {
		return fmt.Errorf("response round %d did not match expected round %d", responseValidRound, uint64(round))
	}
	realLen := len(assetBalancesResponse)
	if realLen != length {
		return fmt.Errorf("response length %d did not match expected length %d", realLen, length)
	}
	scrutinizedElement := assetTransactionsResponse[idx]
	if scrutinizedElement.Sender != sender {
		return fmt.Errorf("response sender %s did not match expected sender %s", scrutinizedElement.Sender, sender)
	}
	return nil
}

var lookupAccountTransactionsResponse []models.Transaction

func weMakeAnyLookupAccountTransactionsCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, lookupAccountTransactionsResponse, globalErrForExamination = indexerClient.NewLookupAccountTransactionsService("").Do(context.Background())
	return nil
}

func theParsedLookupAccountTransactionsResponseShouldBeValidOnRoundAndContainAnArrayOfLenAndElementNumberShouldHaveSender(round, length, idx int, sender string) error {
	if responseValidRound != uint64(round) {
		return fmt.Errorf("response round %d did not match expected round %d", responseValidRound, uint64(round))
	}
	realLen := len(lookupAccountTransactionsResponse)
	if realLen != length {
		return fmt.Errorf("response length %d did not match expected length %d", realLen, length)
	}
	scrutinizedElement := lookupAccountTransactionsResponse[idx]
	if scrutinizedElement.Sender != sender {
		return fmt.Errorf("response sender %s did not match expected sender %s", scrutinizedElement.Sender, sender)
	}
	return nil
}

var lookupBlockResponse models.Block

func weMakeAnyLookupBlockCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	lookupBlockResponse, globalErrForExamination = indexerClient.NewLookupBlockService(0).Do(context.Background())
	return nil
}

func theParsedLookupBlockResponseShouldHaveProposer(proposer string) error {
	if lookupBlockResponse.Proposer != proposer {
		return fmt.Errorf("response proposer %s did not match expected proposer %s", lookupBlockResponse.Proposer, proposer)
	}
	return nil
}

var lookupAccountByIDResponse models.Account

func weMakeAnyLookupAccountByIDCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, lookupAccountByIDResponse, globalErrForExamination = indexerClient.NewLookupAccountByIDService("").Do(context.Background())
	return nil
}

func theParsedLookupAccountByIDResponseShouldHaveAddress(address string) error {
	if lookupAccountByIDResponse.Address != address {
		return fmt.Errorf("response address %s did not match expected address %s", lookupAccountByIDResponse.Address, address)
	}
	return nil
}

var lookupAssetByIDResponse models.Asset

func weMakeAnyLookupAssetByIDCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, lookupAssetByIDResponse, globalErrForExamination = indexerClient.NewLookupAssetByIDService(0).Do(context.Background())
	return nil
}

func theParsedLookupAssetByIDResponseShouldHaveIndex(index int) error {
	if lookupAssetByIDResponse.Index != uint64(index) {
		return fmt.Errorf("response index %d did not match expected index %d", lookupAssetByIDResponse.Index, uint64(index))
	}
	return nil
}

var searchAccountsResponse []models.Account

func weMakeAnySearchAccountsCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, searchAccountsResponse, globalErrForExamination = indexerClient.NewSearchAccountsService().Do(context.Background())
	return nil
}

func theParsedSearchAccountsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAddress(round, length, idx int, address string) error {
	if responseValidRound != uint64(round) {
		return fmt.Errorf("response round %d did not match expected round %d", responseValidRound, uint64(round))
	}
	realLen := len(searchAccountsResponse)
	if realLen != length {
		return fmt.Errorf("response length %d did not match expected length %d", realLen, length)
	}
	scrutinizedElement := searchAccountsResponse[idx]
	if scrutinizedElement.Address != address {
		return fmt.Errorf("response address %s did not match expected address %s", scrutinizedElement.Address, address)
	}
	return nil
}

var searchTransactionsResponse []models.Transaction

func weMakeAnySearchForTransactionsCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, searchTransactionsResponse, globalErrForExamination = indexerClient.NewSearchForTransactionsService().Do(context.Background())
	return nil
}

func theParsedSearchForTransactionsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveSender(round, length, idx int, sender string) error {
	if responseValidRound != uint64(round) {
		return fmt.Errorf("response round %d did not match expected round %d", responseValidRound, uint64(round))
	}
	realLen := len(searchTransactionsResponse)
	if realLen != length {
		return fmt.Errorf("response length %d did not match expected length %d", realLen, length)
	}
	scrutinizedElement := searchTransactionsResponse[idx]
	if scrutinizedElement.Sender != sender {
		return fmt.Errorf("response sender %s did not match expected sender %s", scrutinizedElement.Sender, sender)
	}
	return nil
}

var searchAssetsResponse []models.Asset

func weMakeAnySearchForAssetsCall() error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	responseValidRound, searchAssetsResponse, globalErrForExamination = indexerClient.NewSearchForAssetsService().Do(context.Background())
	return nil
}

func theParsedSearchForAssetsResponseShouldBeValidOnRoundAndTheArrayShouldBeOfLenAndTheElementAtIndexShouldHaveAssetIndex(round, length, idx, expectedIndex int) error {
	if responseValidRound != uint64(round) {
		return fmt.Errorf("response round %d did not match expected round %d", responseValidRound, uint64(round))
	}
	realLen := len(searchAssetsResponse)
	if realLen != length {
		return fmt.Errorf("response length %d did not match expected length %d", realLen, length)
	}
	scrutinizedElement := searchAssetsResponse[idx]
	if scrutinizedElement.Index != uint64(expectedIndex) {
		return fmt.Errorf("response asset index %d did not match expected index %d", scrutinizedElement.Index, uint64(expectedIndex))
	}
	return nil
}

func weMakeALookupAssetBalancesCallAgainstAssetIndexWithLimitLimitAfterAddressRoundCurrencyGreaterThanCurrencyLessThan(index, limit int, afterAddress string, round, currencyGreater, currencyLesser int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = indexerClient.NewLookupAssetBalancesService(uint64(index)).Limit(uint64(limit)).AfterAddress(afterAddress).Round(uint64(round)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).Do(context.Background())
	return nil
}

func weMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseTo(assetIndex int, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit, beforeTime, afterTime, currencyGreater, currencyLesser int, address, addressRole, excludeCloseTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	var excludeCloseToBool bool
	if excludeCloseTo == "false" {
		excludeCloseToBool = false
	} else if excludeCloseTo == "true" {
		excludeCloseToBool = true
	} else {
		return fmt.Errorf("unrecognized expected excludeCloseToBool value %s, allowed values are \"true\" \"false\"", excludeCloseTo)
	}
	_, _, globalErrForExamination = indexerClient.NewLookupAssetTransactionsService(uint64(assetIndex)).NotePrefix([]byte(notePrefix)).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTime(time.Unix(int64(beforeTime), 0)).AfterTime(time.Unix(int64(afterTime), 0)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).Address(address).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).Do(context.Background())
	return nil
}

func weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseTo(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit, beforeTime, afterTime, currencyGreater, currencyLesser, assetIndex int, addressRole, excludeCloseTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	var excludeCloseToBool bool
	if excludeCloseTo == "false" {
		excludeCloseToBool = false
	} else if excludeCloseTo == "true" {
		excludeCloseToBool = true
	} else {
		return fmt.Errorf("unrecognized expected excludeCloseToBool value %s, allowed values are \"true\" \"false\"", excludeCloseTo)
	}
	_, _, globalErrForExamination = indexerClient.NewLookupAccountTransactionsService(account).NotePrefix([]byte(notePrefix)).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTime(time.Unix(int64(beforeTime), 0)).AfterTime(time.Unix(int64(afterTime), 0)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).Do(context.Background())
	return nil
}

func weMakeALookupBlockCallAgainstRound(round int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.NewLookupBlockService(uint64(round)).Do(context.Background())
	return nil
}

func weMakeALookupAccountByIDCallAgainstAccountWithRound(account string, round int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = indexerClient.NewLookupAccountByIDService(account).Round(uint64(round)).Do(context.Background())
	return nil
}

func weMakeALookupAssetByIDCallAgainstAssetIndex(assetIndex int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = indexerClient.NewLookupAssetByIDService(uint64(assetIndex)).Do(context.Background())
	return nil
}

func weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanAndAfterAddress(assetIndex, limit, currencyGreater, currencyLesser int, afterAddress string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = indexerClient.NewSearchAccountsService().AssetID(uint64(assetIndex)).Limit(uint64(limit)).CurrencyLessThan(uint64(currencyLesser)).CurrencyGreaterThan(uint64(currencyGreater)).AfterAddress(afterAddress).Do(context.Background())
	return nil
}

func weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseTo(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit, beforeTime, afterTime, currencyGreater, currencyLesser, assetIndex int, addressRole, excludeCloseTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	var excludeCloseToBool bool
	if excludeCloseTo == "false" {
		excludeCloseToBool = false
	} else if excludeCloseTo == "true" {
		excludeCloseToBool = true
	} else {
		return fmt.Errorf("unrecognized expected excludeCloseToBool value %s, allowed values are \"true\" \"false\"", excludeCloseTo)
	}
	_, _, globalErrForExamination = indexerClient.NewSearchForTransactionsService().Address(account).NotePrefix([]byte(notePrefix)).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTime(time.Unix(int64(beforeTime), 0)).AfterTime(time.Unix(int64(afterTime), 0)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).Do(context.Background())
	return nil
}

func weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndexAndAfterAsset(limit int, creator, name, unit string, assetIndex, afterAsset int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = indexerClient.NewSearchForAssetsService().AssetID(uint64(assetIndex)).Limit(uint64(limit)).Creator(creator).Name(name).Unit(unit).AfterAsset(uint64(afterAsset)).Do(context.Background())
	return nil
}
