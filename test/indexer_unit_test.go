package test

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
)

func IndexerUnitTestContext(s *godog.ScenarioContext) {
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
	s.Step(`^we make any SearchForBlockHeaders call$`, weMakeAnySearchForBlockHeadersCall)
	s.Step(`^the parsed SearchForBlockHeaders response should have a block array of len (\d+) and the element at index (\d+) should have round "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any SearchForAssets call$`, weMakeAnySearchForAssetsCall)
	s.Step(`^the parsed SearchForAssets response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have asset index (\d+)$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make a Lookup Asset Balances call against asset index (\d+) with limit (\d+) afterAddress "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+)$`, weMakeALookupAssetBalancesCallAgainstAssetIndexWithLimitLimitAfterAddressCurrencyGreaterThanCurrencyLessThan)
	s.Step(`^we make a Lookup Asset Transactions call against asset index (\d+) with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) address "([^"]*)" addressRole "([^"]*)" ExcluseCloseTo "([^"]*)"$`, deprecatedWeMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseTo)
	s.Step(`^we make a Search Accounts call with assetID (\d+) limit (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) and round (\d+)$`, weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanAndRound)
	s.Step(`^we make a Lookup Account Transactions call against account "([^"]*)" with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+)$`, weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndex)
	s.Step(`^we make a Lookup Block call against round (\d+)$`, weMakeALookupBlockCallAgainstRound)
	s.Step(`^we make a Lookup Account by ID call against account "([^"]*)" with round (\d+)$`, weMakeALookupAccountByIDCallAgainstAccountWithRound)
	s.Step(`^we make a Lookup Asset by ID call against asset index (\d+)$`, weMakeALookupAssetByIDCallAgainstAssetIndex)
	s.Step(`^mock server recording request paths`, mockServerRecordingRequestPaths)
	s.Step(`^we make a Search For Transactions call with account "([^"]*)" NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) addressRole "([^"]*)" ExcluseCloseTo "([^"]*)" groupid "([^"]*)"$`, weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseToGroupId)
	s.Step(`^we make a Search For BlockHeaders call with minRound (\d+) maxRound (\d+) limit (\d+) nextToken "([^"]*)" beforeTime "([^"]*)" afterTime "([^"]*)" proposers "([^"]*)" expired "([^"]*)" absent "([^"]*)"$`, weMakeASearchForBlockHeadersCallWithMinRoundMaxRoundLimitNextTokenBeforeTimeAfterTimeProposersExpiredAbsent)
	s.Step(`^we make a SearchForAssets call with limit (\d+) creator "([^"]*)" name "([^"]*)" unit "([^"]*)" index (\d+)$`, weMakeASearchForAssetsCallWithLimitCreatorNameUnitIndex)
	s.Step(`^we make a Lookup Account Transactions call against account "([^"]*)" with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) rekeyTo "([^"]*)"$`, weMakeALookupAccountTransactionsCallAgainstAccountWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexRekeyTo)
	s.Step(`^we make a Search Accounts call with assetID (\d+) limit (\d+) currencyGreaterThan (\d+) currencyLessThan (\d+) round (\d+) and authenticating address "([^"]*)"$`, weMakeASearchAccountsCallWithAssetIDLimitCurrencyGreaterThanCurrencyLessThanRoundAndAuthenticatingAddress)
	s.Step(`^we make a Search Accounts call with onlineOnly "([^"]*)"$`, weMakeASearchAccountsCallWithOnlineOnly)
	s.Step(`^we make a Search For Transactions call with account "([^"]*)" NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) assetIndex (\d+) addressRole "([^"]*)" ExcluseCloseTo "([^"]*)" groupid "([^"]*)" rekeyTo "([^"]*)"$`, weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseToGroupidRekeyTo)
	s.Step(`^the parsed SearchAccounts response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have authorizing address "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^the parsed SearchForTransactions response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have rekey-to "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make a Lookup Asset Transactions call against asset index (\d+) with NotePrefix "([^"]*)" TxType "([^"]*)" SigType "([^"]*)" txid "([^"]*)" round (\d+) minRound (\d+) maxRound (\d+) limit (\d+) beforeTime "([^"]*)" afterTime "([^"]*)" currencyGreaterThan (\d+) currencyLessThan (\d+) address "([^"]*)" addressRole "([^"]*)" ExcluseCloseTo "([^"]*)" RekeyTo "([^"]*)"$`, weMakeALookupAssetTransactionsCallAgainstAssetIndexWithNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAddressAddressRoleExcluseCloseToRekeyTo)
	s.Step(`^we make a LookupApplicationLogsByID call with applicationID (\d+) limit (\d+) minRound (\d+) maxRound (\d+) nextToken "([^"]*)" sender "([^"]*)" and txID "([^"]*)"$`, weMakeALookupApplicationLogsByIDCallWithApplicationIDLimitMinRoundMaxRoundNextTokenSenderAndTxID)
	s.Step(`^we make a LookupAccountAssets call with accountID "([^"]*)" assetID (\d+) includeAll "([^"]*)" limit (\d+) next "([^"]*)"$`, weMakeALookupAccountAssetsCallWithAccountIDAssetIDIncludeAllLimitNext)
	s.Step(`^we make a LookupAccountCreatedAssets call with accountID "([^"]*)" assetID (\d+) includeAll "([^"]*)" limit (\d+) next "([^"]*)"$`, weMakeALookupAccountCreatedAssetsCallWithAccountIDAssetIDIncludeAllLimitNext)
	s.Step(`^we make a LookupAccountAppLocalStates call with accountID "([^"]*)" applicationID (\d+) includeAll "([^"]*)" limit (\d+) next "([^"]*)"$`, weMakeALookupAccountAppLocalStatesCallWithAccountIDApplicationIDIncludeAllLimitNext)
	s.Step(`^we make a LookupAccountCreatedApplications call with accountID "([^"]*)" applicationID (\d+) includeAll "([^"]*)" limit (\d+) next "([^"]*)"$`, weMakeALookupAccountCreatedApplicationsCallWithAccountIDApplicationIDIncludeAllLimitNext)
	s.Step(`^we make a Search Accounts call with exclude "([^"]*)"$`, weMakeASearchAccountsCallWithExclude)
	s.Step(`^we make a Lookup Account by ID call against account "([^"]*)" with exclude "([^"]*)"$`, weMakeALookupAccountByIDCallAgainstAccountWithExclude)
	s.Step(`^we make a SearchForApplications call with creator "([^"]*)"$`, weMakeASearchForApplicationsCallWithCreator)
	s.Step(`^we make a Lookup Block call against round (\d+) and header "([^"]*)"$`, weMakeALookupBlockCallAgainstRoundAndHeader)
	s.Step(`^we make a LookupApplicationBoxByIDandName call with applicationID (\d+) with encoded box name "([^"]*)"$`, weMakeALookupApplicationBoxByIDandName)
	s.Step(`^we make a SearchForApplicationBoxes call with applicationID (\d+) with max (\d+) nextToken "([^"]*)"$`, weMakeASearchForApplicationBoxes)
	s.Step(`^the parsed SearchForTransactions response should be valid on round (\d+) and the array should be of len (\d+) and the element at index (\d+) should have hbaddress "([^"]*)"$`, theParsedSearchForTransactionsResponseHasHeartbeatAddress)
	s.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		globalErrForExamination = nil
		return ctx, nil
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

func weMakeAnySearchForBlockHeadersCall() error {
	return weMakeAnyCallTo("indexer", "searchForBlockHeaders")
}

func weMakeAnySearchForAssetsCall() error {
	return weMakeAnyCallTo("indexer", "searchForAssets")
}

func weMakeALookupAssetBalancesCallAgainstAssetIndexWithLimitLimitAfterAddressCurrencyGreaterThanCurrencyLessThan(index, limit int, _ string, currencyGreater, currencyLesser int) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.LookupAssetBalances(uint64(index)).Limit(uint64(limit)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).Do(context.Background())
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
	_, globalErrForExamination = indexerClient.LookupAssetTransactions(uint64(assetIndex)).
		NotePrefix(notePrefixBytes).
		TxType(txType).
		SigType(sigType).
		TXID(txid).
		Round(uint64(round)).
		MinRound(uint64(minRound)).
		MaxRound(uint64(maxRound)).
		Limit(uint64(limit)).
		BeforeTimeString(beforeTime).
		AfterTimeString(afterTime).
		CurrencyGreaterThan(uint64(currencyGreater)).
		CurrencyLessThan(uint64(currencyLesser)).
		AddressString(address).
		AddressRole(addressRole).
		ExcludeCloseTo(excludeCloseToBool).
		RekeyTo(rekeyToBool).
		Do(context.Background())
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

func weMakeASearchAccountsCallWithOnlineOnly(onlineOnly string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.SearchAccounts().OnlineOnly(onlineOnly == "true").Do(context.Background())
	return nil
}

func weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseToGroupId(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser, assetIndex int, addressRole, excludeCloseTo, groupid string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	excludeCloseToBool := excludeCloseTo == "true"
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	groupidBytes, err := base64.StdEncoding.DecodeString(groupid)
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.SearchForTransactions().AddressString(account).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).GroupID(groupidBytes).Do(context.Background())
	return nil
}

func weMakeASearchForBlockHeadersCallWithMinRoundMaxRoundLimitNextTokenBeforeTimeAfterTimeProposersExpiredAbsent(minRound, maxRound, limit int, nextToken, beforeTime, afterTime, proposers, expired, absent string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}

	sfbhBuilder := indexerClient.SearchForBlockHeaders().MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).Next(nextToken).BeforeTimeString(beforeTime).AfterTimeString(afterTime)
	if len(proposers) > 0 {
		proposersArray := strings.Split(proposers, ",")
		sfbhBuilder.Proposers(proposersArray)
	}
	if len(expired) > 0 {
		expiredArray := strings.Split(expired, ",")
		sfbhBuilder.Expired(expiredArray)
	}
	if len(absent) > 0 {
		absentArray := strings.Split(absent, ",")
		sfbhBuilder.Absent(absentArray)
	}
	_, globalErrForExamination = sfbhBuilder.Do(context.Background())
	return nil
}

func weMakeASearchForTransactionsCallWithAccountNotePrefixTxTypeSigTypeTxidRoundMinRoundMaxRoundLimitBeforeTimeAfterTimeCurrencyGreaterThanCurrencyLessThanAssetIndexAddressRoleExcluseCloseToGroupidRekeyTo(account, notePrefix, txType, sigType, txid string, round, minRound, maxRound, limit int, beforeTime, afterTime string, currencyGreater, currencyLesser, assetIndex int, addressRole, excludeCloseTo, groupid, rekeyTo string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	excludeCloseToBool := excludeCloseTo == "true"
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	groupidBytes, err := base64.StdEncoding.DecodeString(groupid)
	if err != nil {
		return err
	}
	rekeyToBool := rekeyTo == "true"
	_, globalErrForExamination = indexerClient.SearchForTransactions().AddressString(account).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Limit(uint64(limit)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AssetID(uint64(assetIndex)).AddressRole(addressRole).ExcludeCloseTo(excludeCloseToBool).GroupID(groupidBytes).RekeyTo(rekeyToBool).Do(context.Background())
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
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = indexerClient.SearchForAssets().AssetID(uint64(index)).Limit(uint64(limit)).Creator(creator).Name(name).Unit(unit).Do(context.Background())
	return nil
}

func weMakeALookupApplicationLogsByIDCallWithApplicationIDLimitMinRoundMaxRoundNextTokenSenderAndTxID(appID, limit, minRound, maxRound int, nextToken, sender, txID string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	indexerClient.LookupApplicationLogsByID(uint64(appID)).Limit(uint64(limit)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).Next(nextToken).SenderAddress(sender).Txid(txID).Do(context.Background())
	return nil
}

func parseBool(s string) (bool, error) {
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	return false, fmt.Errorf("parseBool() cannot parse \"%s\"", s)
}

func weMakeALookupAccountAssetsCallWithAccountIDAssetIDIncludeAllLimitNext(accountID string, assetID int, includeAll string, limit int, next string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	includeAllBool, err := parseBool(includeAll)
	if err != nil {
		return err
	}
	indexerClient.LookupAccountAssets(accountID).AssetID(uint64(assetID)).IncludeAll(includeAllBool).Limit(uint64(limit)).Next(next).Do(context.Background())
	return nil
}

func weMakeALookupAccountCreatedAssetsCallWithAccountIDAssetIDIncludeAllLimitNext(accountID string, assetID int, includeAll string, limit int, next string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	includeAllBool, err := parseBool(includeAll)
	if err != nil {
		return err
	}
	indexerClient.LookupAccountCreatedAssets(accountID).AssetID(uint64(assetID)).IncludeAll(includeAllBool).Limit(uint64(limit)).Next(next).Do(context.Background())
	return nil
}

func weMakeALookupAccountAppLocalStatesCallWithAccountIDApplicationIDIncludeAllLimitNext(accountID string, appID int, includeAll string, limit int, next string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	includeAllBool, err := parseBool(includeAll)
	if err != nil {
		return err
	}
	indexerClient.LookupAccountAppLocalStates(accountID).ApplicationID(uint64(appID)).IncludeAll(includeAllBool).Limit(uint64(limit)).Next(next).Do(context.Background())
	return nil
}

func weMakeALookupAccountCreatedApplicationsCallWithAccountIDApplicationIDIncludeAllLimitNext(accountID string, appID int, includeAll string, limit int, next string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	includeAllBool, err := parseBool(includeAll)
	if err != nil {
		return err
	}
	indexerClient.LookupAccountCreatedApplications(accountID).ApplicationID(uint64(appID)).IncludeAll(includeAllBool).Limit(uint64(limit)).Next(next).Do(context.Background())
	return nil
}

func weMakeASearchAccountsCallWithExclude(exclude string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	indexerClient.SearchAccounts().Exclude(strings.Split(exclude, ",")).Do(context.Background())
	return nil
}

func weMakeALookupAccountByIDCallAgainstAccountWithExclude(account, exclude string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	indexerClient.LookupAccountByID(account).Exclude(strings.Split(exclude, ",")).Do(context.Background())
	return nil
}

func weMakeASearchForApplicationsCallWithCreator(creator string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	indexerClient.SearchForApplications().Creator(creator).Do(context.Background())
	return nil
}

func weMakeALookupBlockCallAgainstRoundAndHeader(round int, headerOnly string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}

	headerOnlyBool, err := parseBool(headerOnly)
	if err != nil {
		return weMakeALookupBlockCallAgainstRound(round)
	}

	_, globalErrForExamination = indexerClient.LookupBlock(uint64(round)).HeaderOnly(headerOnlyBool).Do(context.Background())
	return nil
}

func weMakeALookupApplicationBoxByIDandName(appId int, encodedBoxName string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	decodedBoxNames, err := parseAppArgs(encodedBoxName)
	if err != nil {
		return err
	}
	indexerClient.LookupApplicationBoxByIDAndName(uint64(appId), decodedBoxNames[0]).Do(context.Background())
	return nil
}

func weMakeASearchForApplicationBoxes(appId int, limit int, next string) error {
	indexerClient, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	indexerClient.SearchForApplicationBoxes(uint64(appId)).Limit(uint64(limit)).Next(next).Do(context.Background())
	return nil
}

func theParsedSearchForTransactionsResponseHasHeartbeatAddress(round int, length int, index int, hbAddress string) error {
	resp := response.(models.TransactionsResponse)

	if resp.CurrentRound != uint64(round) {
		return fmt.Errorf("Expected round %d, got %d", round, resp.CurrentRound)
	}
	if len(resp.Transactions) != length {
		return fmt.Errorf("Expected %d transactions, got %d", length, len(resp.Transactions))
	}

	hbTxn := resp.Transactions[index]
	if hbTxn.HeartbeatTransaction.HbAddress != hbAddress {
		return fmt.Errorf("Expected heartbeat address %s, got %s", hbAddress, hbTxn.HeartbeatTransaction.HbAddress)
	}

	return nil
}
