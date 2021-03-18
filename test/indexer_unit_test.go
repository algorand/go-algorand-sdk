package test

import (
	"context"
	"encoding/base64"

	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/cucumber/godog"
)

func IndexerUnitTestContext(s *godog.Suite) {
	s.Step(`^we make any LookupAssetBalances call$`, weMakeAnyLookupAssetBalancesCall)
	s.Step(`^the parsed LookupAssetBalances response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have address "([^"]*)" amount (\d+) and frozen state "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any LookupAssetTransactions call$`, weMakeAnyLookupAssetTransactionsCall)
	s.Step(`^the parsed LookupAssetTransactions response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any LookupAccountTransactions call$`, weMakeAnyLookupAccountTransactionsCall)
	s.Step(`^the parsed LookupAccountTransactions response should be valid on round (\d+), and contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any LookupBlock call$`, weMakeAnyLookupBlockCall)
	s.Step(`^the parsed LookupBlock response should have previous block hash "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any LookupAccountByID call$`, weMakeAnyLookupAccountByIDCall)
	s.Step(`^the parsed LookupAccountByID response should have address "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any LookupAssetByID call$`, weMakeAnyLookupAssetByIDCall)
	s.Step(`^the parsed LookupAssetByID response should have index (\d+)$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any SearchAccounts call$`, weMakeAnySearchAccountsCall)
	s.Step(`^the parsed SearchAccounts response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have address "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any SearchForTransactions call$`, weMakeAnySearchForTransactionsCall)
	s.Step(`^the parsed SearchForTransactions response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any SearchForAssets call$`, weMakeAnySearchForAssetsCall)
	s.Step(`^the parsed SearchForAssets response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have asset index (\d+)$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make a Lookup Asset Balances call against asset index (\d+) with limit (\d+) afterAddress "([^"]*)" round (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+)$`, weMakeALookupAssetBalancesCallAgainstAssetIndexWithLimitLimitAfterAddressRoundCurrencyGreaterThanCurrencyLessThan)
	s.Step(`^we make a Lookup Asset Transactions call against asset index (\d+) with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) address "([^"]*)" addressRole "([^"]*)" ExcluseCloseTo "([^"]*)"$`, deprecatedWeMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseTo)
	s.Step(`^we make a Search Accounts call with assetID (\d+) limit (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) and round (\d+)$`, weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanAndRound)
	s.Step(`^we make a Lookup Account Transactions call against account "([^"]*)" with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+)$`, weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndex)
	s.Step(`^we make a Lookup Block call against round (\d+)$`, weMakeALookupBlockCallAgainstRound)
	s.Step(`^we make a Lookup Account by ID call against account "([^"]*)" with round (\d+)$`, weMakeALookupAccountByIDCallAgainstAccountWithRound)
	s.Step(`^we make a Lookup Asset by ID call against asset index (\d+)$`, weMakeALookupAssetByIDCallAgainstAssetIndex)
	s.Step(`^we make a Search For Transactions call with account "([^"]*)" NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime (\d+) afterTime (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) addressRole "([^"]*)" ExcluseCloseTo "([^"]*)"$`, weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseTo)
	s.Step(`^we make a SearchForAssets call with limit (\d+) creator "([^"]*)" name "([^"]*)" unit "([^"]*)" index (\d+) and afterAsset (\d+)$`, weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndexAndAfterAsset)
	s.Step(`^mock server recording request paths`, mockServerRecordingRequestPaths)
	s.Step(`^we make a Search For Transactions call with account "([^"]*)" NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) addressRole "([^"]*)" ExcluseCloseTo "([^"]*)"$`, weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseTo)
	s.Step(`^we make a SearchForAssets call with limit (\d+) creator "([^"]*)" name "([^"]*)" unit "([^"]*)" index (\d+)$`, weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndex)
	s.Step(`^we make a Lookup Account Transactions call against account "([^"]*)" with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) rekeyTo "([^"]*)"$`, weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexRekeyTo)
	s.Step(`^we make a Search Accounts call with assetID (\d+) limit (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) round (\d+) and authenticating address "([^"]*)"$`, weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanRoundAndAuthenticatingAddress)
	s.Step(`^we make a Search For Transactions call with account "([^"]*)" NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) addressRole "([^"]*)" ExcluseCloseTo "([^"]*)" rekeyTo "([^"]*)"$`, weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseToRekeyTo)
	s.Step(`^the parsed SearchAccounts response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have authorizing address "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^the parsed SearchForTransactions response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have rekey-to "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make a Lookup Asset Transactions call against asset index (\d+) with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) address "([^"]*)" addressRole "([^"]*)" ExcluseCloseTo "([^"]*)" RekeyTo "([^"]*)"$`, weMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseToRekeyTo)
	s.BeforeScenario(func(interface{}) {
		globalErrForExamination = nil
	})
}

func weMakeAnyLookupAssetBalancesCall() error {
	return weMakeAnyCallTo("indexer", "lookupAssetBalances")
}

func weMakeAnyLookupAssetTransactionsCall() error {
	return weMakeAnyCallTo("indexer", "lookupAssetTransactions")
}

func weMakeAnyLookupAccountTransactionsCall() error {
	return weMakeAnyCallTo("indexer", "lookupAccountTransactions")
}

func weMakeAnyLookupBlockCall() error {
	return weMakeAnyCallTo("indexer", "lookupBlock")
}

func weMakeAnyLookupAccountByIDCall() error {
	return weMakeAnyCallTo("indexer", "lookupAccountByID")
}

func weMakeAnyLookupAssetByIDCall() error {
	return weMakeAnyCallTo("indexer", "lookupAssetByID")
}

func weMakeAnySearchAccountsCall() error {
	return weMakeAnyCallTo("indexer", "searchForAccounts")
}

func weMakeAnySearchForTransactionsCall() error {
	return weMakeAnyCallTo("indexer", "searchForTransactions")
}

func weMakeAnySearchForAssetsCall() error {
	return weMakeAnyCallTo("indexer", "searchForAssets")
}

func weMakeALookupAssetBalancesCallAgainstAssetIndexWithLimitLimitAfterAddressRoundCurrencyGreaterThanCurrencyLessThan(index, limit int, _ string, round, currencyGreater, currencyLesser int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.LookupAssetBalances(uint64(index)).Limit(uint64(limit)).Round(uint64(round)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).Do(context.Background())
	return nil
}

func deprecatedWeMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseTo(assetIndex int, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser int, address, addressRole, excludeCloseTo string) error {
	return weMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseToRekeyTo(assetIndex, notePrefix, txType, sigType, txid, round, minRound, maxRound, limit, beforeTime, afterTime, currencyGreater, currencyLesser, address, addressRole, excludeCloseTo, "")
}

func weMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseToRekeyTo(assetIndex int, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser int, address, addressRole, excludeCloseTo, rekeyTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	excludeCloseToBool := excludeCloseTo == "true"
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	rekeyToBool := rekeyTo == "true"
	_, globalErrForExamination = indexerClient.LookupAssetTransactions(uint64(assetIndex)).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AddressString(address).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).RekeyTo(rekeyToBool).Do(context.Background())
	return nil
}

func weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndex(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser, assetIndex int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.LookupAccountTransactions(account).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).Do(context.Background())
	return nil
}

func weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexRekeyTo(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser, assetIndex int, rekeyTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	rekeyToBool := rekeyTo == "true"
	_, globalErrForExamination = indexerClient.LookupAccountTransactions(account).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).RekeyTo(rekeyToBool).Do(context.Background())
	return nil
}

func weMakeALookupBlockCallAgainstRound(round int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.LookupBlock(uint64(round)).Do(context.Background())
	return nil
}

func weMakeALookupAccountByIDCallAgainstAccountWithRound(account string, round int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = indexerClient.LookupAccountByID(account).Round(uint64(round)).Do(context.Background())
	return nil
}

func weMakeALookupAssetByIDCallAgainstAssetIndex(assetIndex int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = indexerClient.LookupAssetByID(uint64(assetIndex)).Do(context.Background())
	return nil
}

func weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanRoundAndAuthenticatingAddress(assetIndex, limit, currencyGreater, currencyLesser, round int, authAddress string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.SearchAccounts().
		AssetID(uint64(assetIndex)).
		Limit(uint64(limit)).
		CurrencyLessThan(uint64(currencyLesser)).
		CurrencyGreaterThan(uint64(currencyGreater)).
		Round(uint64(round)).
		AuthAddress(authAddress).
		Do(context.Background())
	return nil
}

func weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseTo(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser, assetIndex int, addressRole, excludeCloseTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	excludeCloseToBool := excludeCloseTo == "true"
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.SearchForTransactions().AddressString(account).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).Do(context.Background())
	return nil
}

func weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseToRekeyTo(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser, assetIndex int, addressRole, excludeCloseTo, rekeyTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	excludeCloseToBool := excludeCloseTo == "true"
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	rekeyToBool := rekeyTo == "true"
	_, globalErrForExamination = indexerClient.SearchForTransactions().AddressString(account).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).RekeyTo(rekeyToBool).Do(context.Background())
	return nil
}

func weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndexAndAfterAsset(limit int, creator, name, unit string, assetIndex, _ int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.SearchForAssets().AssetID(uint64(assetIndex)).Limit(uint64(limit)).Creator(creator).Name(name).Unit(unit).Do(context.Background())
	return nil
}

func weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanAndRound(assetID, limit, currencyGreater, currencyLesser, round int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.SearchAccounts().AssetID(uint64(assetID)).Limit(uint64(limit)).CurrencyLessThan(uint64(currencyLesser)).CurrencyGreaterThan(uint64(currencyGreater)).Round(uint64(round)).Do(context.Background())
	return nil
}

func weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndex(limit int, creator, name, unit string, index int) error {
	return weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndexAndAfterAsset(limit, creator, name, unit, index, 0)
}
