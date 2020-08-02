package test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/encoding/json"
)

func IndexerIntegrationTestContext(s *godog.Suite) {
	s.Step(`^indexer client (\d+) at "([^"]*)" port (\d+) with token "([^"]*)"$`, indexerClientAtPortWithToken)
	s.Step(`^I use (\d+) to check the services health$`, iUseToCheckTheServicesHealth)
	s.Step(`^I receive status code (\d+)$`, iReceiveStatusCode)
	s.Step(`^I use (\d+) to lookup block (\d+)$`, iUseToLookupBlock)
	s.Step(`^The block was confirmed at (\d+), contains (\d+) transactions, has the previous block hash "([^"]*)"$`, theBlockWasConfirmedAtContainsTransactionsHasThePreviousBlockHash)
	s.Step(`^I use (\d+) to lookup account "([^"]*)" at round (\d+)$`, iUseToLookupAccountAtRound)
	s.Step(`^The account has (\d+) assets, the first is asset (\d+) has a frozen status of "([^"]*)" and amount (\d+)\.$`, theAccountHasAssetsTheFirstIsAssetHasAFrozenStatusOfAndAmount)
	s.Step(`^The account created (\d+) assets, the first is asset (\d+) is named "([^"]*)" with a total amount of (\d+) "([^"]*)"$`, theAccountCreatedAssetsTheFirstIsAssetIsNamedWithATotalAmountOf)
	s.Step(`^The account has (\d+) μalgos and (\d+) assets, (\d+) has (\d+)$`, theAccountHasAlgosAndAssetsHas)
	s.Step(`^I use (\d+) to lookup asset (\d+)$`, iUseToLookupAsset)
	s.Step(`^The asset found has: "([^"]*)", "([^"]*)", "([^"]*)", (\d+), "([^"]*)", (\d+), "([^"]*)"$`, theAssetFoundHas)
	s.Step(`^I use (\d+) to lookup asset balances for (\d+) with (\d+), (\d+), (\d+) and token "([^"]*)"$`, iUseToLookupAssetBalancesForWithAndToken)
	s.Step(`^There are (\d+) with the asset, the first is "([^"]*)" has "([^"]*)" and (\d+)$`, thereAreWithTheAssetTheFirstIsHasAnd)
	s.Step(`^I get the next page using (\d+) to lookup asset balances for (\d+) with (\d+), (\d+), (\d+)$`, iGetTheNextPageUsingToLookupAssetBalancesForWith)
	s.Step(`^I use (\d+) to search for an account with (\d+), (\d+), (\d+), (\d+) and token "([^"]*)"$`, iUseToSearchForAnAccountWithAndToken)
	s.Step(`^There are (\d+), the first has (\d+), (\d+), (\d+), (\d+), "([^"]*)", (\d+), "([^"]*)", "([^"]*)"$`, thereAreTheFirstHas)
	s.Step(`^The first account is online and has "([^"]*)", (\d+), (\d+), (\d+), "([^"]*)", "([^"]*)"$`, theFirstAccountIsOnlineAndHas)
	s.Step(`^I get the next page using (\d+) to search for an account with (\d+), (\d+), (\d+) and (\d+)$`, iGetTheNextPageUsingToSearchForAnAccountWithAnd)
	s.Step(`^I use (\d+) to search for transactions with (\d+), "([^"]*)", "([^"]*)", "([^"]*)", "([^"]*)", (\d+), (\d+), (\d+), (\d+), "([^"]*)", "([^"]*)", (\d+), (\d+), "([^"]*)", "([^"]*)", "([^"]*)" and token "([^"]*)"$`, iUseToSearchForTransactionsWithAndToken)
	s.Step(`^there are (\d+) transactions in the response, the first is "([^"]*)"\.$`, thereAreTransactionsInTheResponseTheFirstIs)
	s.Step(`^Every transaction has tx-type "([^"]*)"$`, everyTransactionHasTxtype)
	s.Step(`^Every transaction has sig-type "([^"]*)"$`, everyTransactionHasSigtype)
	s.Step(`^Every transaction has round (\d+)$`, everyTransactionHasRoundEqual)
	s.Step(`^Every transaction has round >= (\d+)$`, everyTransactionHasRoundGreaterThan)
	s.Step(`^Every transaction has round <= (\d+)$`, everyTransactionHasRoundLessThan)
	s.Step(`^Every transaction works with asset-id (\d+)$`, everyTransactionWorksWithAssetid)
	s.Step(`^Every transaction is older than "([^"]*)"$`, everyTransactionIsOlderThan)
	s.Step(`^Every transaction is newer than "([^"]*)"$`, everyTransactionIsNewerThan)
	s.Step(`^Every transaction moves between (\d+) and (\d+) currency$`, everyTransactionMovesBetweenAndCurrency)
	s.Step(`^I use (\d+) to search for all "([^"]*)" transactions$`, iUseToSearchForAllTransactions)
	s.Step(`^I use (\d+) to search for all (\d+) asset transactions$`, iUseToSearchForAllAssetTransactions)
	s.Step(`^I get the next page using (\d+) to search for transactions with (\d+) and (\d+)$`, iGetTheNextPageUsingToSearchForTransactionsWithAnd)
	s.Step(`^I use (\d+) to search for assets with (\d+), (\d+), "([^"]*)", "([^"]*)", "([^"]*)", and token "([^"]*)"$`, iUseToSearchForAssetsWithAndToken)
	s.Step(`^there are (\d+) assets in the response, the first is (\d+)\.$`, thereAreAssetsInTheResponseTheFirstIs)

	//@indexer.applications
	s.Step(`^I use (\d+) to search for applications with (\d+), (\d+), and token "([^"]*)"$`, iUseToSearchForApplicationsWithAndToken)
	s.Step(`^the parsed response should equal "([^"]*)"\.$`, theParsedResponseShouldEqual)
	s.Step(`^I use (\d+) to lookup application with (\d+)$`, iUseToLookupApplicationWith)
	s.Step(`^I use (\d+) to search for transactions with (\d+), "([^"]*)", "([^"]*)", "([^"]*)", "([^"]*)", (\d+), (\d+), (\d+), (\d+), "([^"]*)", "([^"]*)", (\d+), (\d+), "([^"]*)", "([^"]*)", "([^"]*)", (\d+) and token "([^"]*)"$`, iUseToSearchForTransactionsWithAppIdAndToken)
	s.Step(`^I use (\d+) to search for an account with (\d+), (\d+), (\d+), (\d+), "([^"]*)", (\d+) and token "([^"]*)"$`, iUseToSearchForAnAccountWithAppIdAndToken)

	s.BeforeScenario(func(interface{}) {
	})
}

var indexerClients = make(map[int]*indexer.Client)

func indexerClientAtPortWithToken(clientNum int, host string, port int, token string) error {
	portAsString := strconv.Itoa(port)
	fullHost := "http://" + host + ":" + portAsString
	ic, err := indexer.MakeClient(fullHost, token)
	indexerClients[clientNum] = ic
	return err
}

var indexerHealthCheckResponse models.HealthCheckResponse
var indexerHealthCheckError error

func iUseToCheckTheServicesHealth(clientNum int) error {
	ic := indexerClients[clientNum]
	indexerHealthCheckResponse, indexerHealthCheckError = ic.HealthCheck().Do(context.Background())
	return nil
}

func iReceiveStatusCode(code int) error {
	if code == 200 && indexerHealthCheckError == nil {
		return nil
	}
	return fmt.Errorf("Did not receive expected error code: %v", indexerHealthCheckError)
}

var indexerBlockResponse models.Block

func iUseToLookupBlock(clientNum, blockNum int) error {
	ic := indexerClients[clientNum]
	var err error
	indexerBlockResponse, err = ic.LookupBlock(uint64(blockNum)).Do(context.Background())
	return err
}

func theBlockWasConfirmedAtContainsTransactionsHasThePreviousBlockHash(timestamp, numTransactions int, prevBlockHash string) error {
	err := comparisonCheck("timestamp", uint64(timestamp), indexerBlockResponse.Timestamp)
	if err != nil {
		return err
	}
	err = comparisonCheck("number of transactions", numTransactions, len(indexerBlockResponse.Transactions))
	if err != nil {
		return err
	}
	err = comparisonCheck("previous blockhash", prevBlockHash, base64.StdEncoding.EncodeToString(indexerBlockResponse.PreviousBlockHash))
	return err
}

var indexerAccountResponse models.Account

func iUseToLookupAccountAtRound(clientNum int, account string, round int) error {
	ic := indexerClients[clientNum]
	var err error
	_, indexerAccountResponse, err = ic.LookupAccountByID(account).Round(uint64(round)).Do(context.Background())
	return err
}

func theAccountHasAssetsTheFirstIsAssetHasAFrozenStatusOfAndAmount(numAssets, firstAssetIndex int, firstAssetFrozenStatus string, firstAssetAmount int) error {
	err := comparisonCheck("number of asset holdings", numAssets, len(indexerAccountResponse.Assets))
	if err != nil {
		return err
	}
	assetUnderScrutiny := indexerAccountResponse.Assets[0]
	err = comparisonCheck("first asset index", uint64(firstAssetIndex), assetUnderScrutiny.AssetId)
	if err != nil {
		return err
	}
	frozen, err := strconv.ParseBool(firstAssetFrozenStatus)
	if err != nil {
		frozen = false
	}
	err = comparisonCheck("first asset frozen status", frozen, assetUnderScrutiny.IsFrozen)
	if err != nil {
		return err
	}
	err = comparisonCheck("first asset amount", uint64(firstAssetAmount), assetUnderScrutiny.Amount)
	return err
}

func theAccountCreatedAssetsTheFirstIsAssetIsNamedWithATotalAmountOf(numCreatedAssets, firstAssetIndex int, assetName string, assetIssuance int, assetUnitName string) error {
	err := comparisonCheck("number of created assets", numCreatedAssets, len(indexerAccountResponse.CreatedAssets))
	if err != nil {
		return err
	}
	assetUnderScrutiny := indexerAccountResponse.CreatedAssets[0]
	err = comparisonCheck("first created asset index", uint64(firstAssetIndex), assetUnderScrutiny.Index)
	if err != nil {
		return err
	}
	err = comparisonCheck("first created asset name", assetName, assetUnderScrutiny.Params.Name)
	if err != nil {
		return err
	}
	err = comparisonCheck("first created asset issuance", uint64(assetIssuance), assetUnderScrutiny.Params.Total)
	if err != nil {
		return err
	}
	err = comparisonCheck("first created asset unit name", assetUnitName, assetUnderScrutiny.Params.UnitName)
	return err
}

func theAccountHasAlgosAndAssetsHas(microAlgos, numAssets, assetIndex, assetAmount int) error {
	err := comparisonCheck("microalgo balance", uint64(microAlgos), indexerAccountResponse.Amount)
	if err != nil {
		return err
	}
	err = comparisonCheck("number of asset holdings", numAssets, len(indexerAccountResponse.Assets))
	if err != nil {
		return err
	}
	if numAssets == 0 || assetIndex == 0 {
		return nil
	}
	assetUnderScrutiny, err := findAssetInHoldingsList(indexerAccountResponse.Assets, uint64(assetIndex))
	if err != nil {
		return err
	}
	err = comparisonCheck(fmt.Sprintf("amount for asset %d", uint64(assetIndex)), uint64(assetAmount), assetUnderScrutiny.Amount)
	return err
}

var indexerAssetResponse models.Asset

func iUseToLookupAsset(clientNum, assetId int) error {
	ic := indexerClients[clientNum]
	var err error
	_, indexerAssetResponse, err = ic.LookupAssetByID(uint64(assetId)).Do(context.Background())
	return err
}

func theAssetFoundHas(assetName, assetUnits, assetCreator string, assetDecimals int, assetDefaultFrozen string, assetIssuance int, assetClawbackAddress string) error {
	err := comparisonCheck("asset name", assetName, indexerAssetResponse.Params.Name)
	if err != nil {
		return err
	}
	err = comparisonCheck("asset units", assetUnits, indexerAssetResponse.Params.UnitName)
	if err != nil {
		return err
	}
	err = comparisonCheck("asset creator", assetCreator, indexerAssetResponse.Params.Creator)
	if err != nil {
		return err
	}
	err = comparisonCheck("asset decimals", uint64(assetDecimals), indexerAssetResponse.Params.Decimals)
	if err != nil {
		return err
	}
	frozen, err := strconv.ParseBool(assetDefaultFrozen)
	if err != nil {
		frozen = false
	}
	err = comparisonCheck("asset default frozen state", frozen, indexerAssetResponse.Params.DefaultFrozen)
	if err != nil {
		return err
	}
	err = comparisonCheck("asset issuance", uint64(assetIssuance), indexerAssetResponse.Params.Total)
	if err != nil {
		return err
	}
	err = comparisonCheck("asset clawback address", assetClawbackAddress, indexerAssetResponse.Params.Clawback)
	return err
}

var indexerAssetBalancesResponse models.AssetBalancesResponse

func iUseToLookupAssetBalancesForWithAndToken(clientNum, assetId, currencyGreater, currencyLesser, limit int, token string) error {
	ic := indexerClients[clientNum]
	var err error
	indexerAssetBalancesResponse, err = ic.LookupAssetBalances(uint64(assetId)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).Limit(uint64(limit)).NextToken(token).Do(context.Background())
	return err
}

func thereAreWithTheAssetTheFirstIsHasAnd(numAccountsWithAsset int, firstHolder, firstHolderIsFrozen string, firstHolderAmount int) error {
	err := comparisonCheck("number of asset holders", numAccountsWithAsset, len(indexerAssetBalancesResponse.Balances))
	if err != nil {
		return err
	}
	if numAccountsWithAsset == 0 {
		return nil
	}
	accountUnderScrutiny := indexerAssetBalancesResponse.Balances[0]
	err = comparisonCheck("first account holder", firstHolder, accountUnderScrutiny.Address)
	if err != nil {
		return err
	}
	frozen, err := strconv.ParseBool(firstHolderIsFrozen)
	if err != nil {
		frozen = false
	}
	err = comparisonCheck("first account holder frozen state", frozen, accountUnderScrutiny.IsFrozen)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account holder asset balance", uint64(firstHolderAmount), accountUnderScrutiny.Amount)
	return err
}

func iGetTheNextPageUsingToLookupAssetBalancesForWith(clientNum, assetId, currencyGreater, currencyLesser, limit int) error {
	ic := indexerClients[clientNum]
	var err error
	indexerAssetBalancesResponse, err = ic.LookupAssetBalances(uint64(assetId)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).Limit(uint64(limit)).NextToken(indexerAssetBalancesResponse.NextToken).Do(context.Background())
	return err
}

var indexerSearchAccountsResponse models.AccountsResponse

func iUseToSearchForAnAccountWithAndToken(clientNum, assetIndex, limit, currencyGreater, currencyLesser int, token string) error {
	ic := indexerClients[clientNum]
	var err error
	indexerSearchAccountsResponse, err = ic.SearchAccounts().AssetID(uint64(assetIndex)).Limit(uint64(limit)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).NextToken(token).Do(context.Background())
	return err
}

func thereAreTheFirstHas(numAccounts, firstAccountPendingRewards, rewardsBase, rewards, withoutRewards int, address string, amount int, accountStatus, accountType string) error {
	err := comparisonCheck("number of found accounts", numAccounts, len(indexerSearchAccountsResponse.Accounts))
	if err != nil {
		return err
	}
	if numAccounts == 0 {
		return nil
	}
	accountUnderScrutiny := indexerSearchAccountsResponse.Accounts[0]
	err = comparisonCheck("first account pending rewards", uint64(firstAccountPendingRewards), accountUnderScrutiny.PendingRewards)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account rewards base", uint64(rewardsBase), accountUnderScrutiny.RewardBase)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account rewards", uint64(rewards), accountUnderScrutiny.Rewards)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account without-rewards", uint64(withoutRewards), accountUnderScrutiny.AmountWithoutPendingRewards)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account address", address, accountUnderScrutiny.Address)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account balance", uint64(amount), accountUnderScrutiny.Amount)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account status", accountStatus, accountUnderScrutiny.Status)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account type", accountType, accountUnderScrutiny.SigType)
	return err
}

func theFirstAccountIsOnlineAndHas(address string, keyDilution, firstValid, lastValid int, voteKey, selKey string) error {
	accountUnderScrutiny := indexerSearchAccountsResponse.Accounts[0]
	err := comparisonCheck("first account online state", "Online", accountUnderScrutiny.Status)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account address", address, accountUnderScrutiny.Address)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account key dilution", uint64(keyDilution), accountUnderScrutiny.Participation.VoteKeyDilution)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account partkey firstvalid", uint64(firstValid), accountUnderScrutiny.Participation.VoteFirstValid)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account partkey lastvalid", uint64(lastValid), accountUnderScrutiny.Participation.VoteLastValid)
	if err != nil {
		return err
	}
	voteKeyString := base64.StdEncoding.EncodeToString(accountUnderScrutiny.Participation.VoteParticipationKey)
	selKeyString := base64.StdEncoding.EncodeToString(accountUnderScrutiny.Participation.SelectionParticipationKey)
	err = comparisonCheck("first account votekey b64", voteKey, voteKeyString)
	if err != nil {
		return err
	}
	err = comparisonCheck("first account selkey b64", selKey, selKeyString)
	return err
}

func iGetTheNextPageUsingToSearchForAnAccountWithAnd(clientNum, assetId, limit, currencyGreater, currencyLesser int) error {
	ic := indexerClients[clientNum]
	var err error
	indexerSearchAccountsResponse, err = ic.SearchAccounts().AssetID(uint64(assetId)).Limit(uint64(limit)).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).NextToken(indexerSearchAccountsResponse.NextToken).Do(context.Background())
	return err
}

var indexerTransactionsResponse models.TransactionsResponse

func iUseToSearchForTransactionsWithAndToken(clientNum, limit int, notePrefix, txType, sigType, txid string, round, minRound, maxRound, assetId int, beforeTime, afterTime string, currencyGreater, currencyLesser int, address, addressRole, excludeCloseTo, token string) error {
	ic := indexerClients[clientNum]
	var err error
	notePrefixBytes, err := base64.StdEncoding.DecodeString(notePrefix)
	if err != nil {
		return err
	}
	excludeBool, err := strconv.ParseBool(excludeCloseTo)
	if err != nil {
		excludeBool = false
	}
	indexerTransactionsResponse, err = ic.SearchForTransactions().Limit(uint64(limit)).NotePrefix(notePrefixBytes).TxType(txType).SigType(sigType).TXID(txid).Round(uint64(round)).MinRound(uint64(minRound)).MaxRound(uint64(maxRound)).AssetID(uint64(assetId)).BeforeTimeString(beforeTime).AfterTimeString(afterTime).CurrencyGreaterThan(uint64(currencyGreater)).CurrencyLessThan(uint64(currencyLesser)).AddressString(address).AddressRole(addressRole).ExcludeCloseTo(excludeBool).NextToken(token).Do(context.Background())
	return err
}

func thereAreTransactionsInTheResponseTheFirstIs(numTransactions int, firstTransactionTxid string) error {
	err := comparisonCheck("number of transactions", numTransactions, len(indexerTransactionsResponse.Transactions))
	if err != nil {
		return err
	}
	if numTransactions == 0 {
		return nil
	}
	txnUnderScrutiny := indexerTransactionsResponse.Transactions[0]
	err = comparisonCheck("first txn txid", firstTransactionTxid, txnUnderScrutiny.Id)
	return err
}

func everyTransactionHasTxtype(txType string) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		err := comparisonCheck(fmt.Sprintf("tx type for returned transaction %d", idx), txType, txn.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func everyTransactionHasSigtype(sigType string) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		err := comparisonCheck(fmt.Sprintf("sig type for returned transaction %d", idx), sigType, getSigtypeFromTransaction(txn))
		if err != nil {
			return err
		}
	}
	return nil
}

func everyTransactionHasRoundEqual(round int) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		err := comparisonCheck(fmt.Sprintf("confirmed round for returned transaction %d", idx), uint64(round), txn.ConfirmedRound)
		if err != nil {
			return err
		}
	}
	return nil
}

func everyTransactionHasRoundGreaterThan(round int) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		if txn.ConfirmedRound < uint64(round) {
			return fmt.Errorf("transaction number %d was confirmed in round %d, but expected a number higher than %d", idx, txn.ConfirmedRound, uint64(round))
		}
	}
	return nil
}

func everyTransactionHasRoundLessThan(round int) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		if txn.ConfirmedRound > uint64(round) {
			return fmt.Errorf("transaction number %d was confirmed in round %d, but expected a number lower than %d", idx, txn.ConfirmedRound, uint64(round))
		}
	}
	return nil
}

func everyTransactionWorksWithAssetid(assetId int) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		actualId := uint64(0)
		if txn.AssetTransferTransaction.AssetId != 0 {
			actualId = txn.AssetTransferTransaction.AssetId
		}
		if txn.AssetFreezeTransaction.AssetId != 0 {
			actualId = txn.AssetFreezeTransaction.AssetId
		}
		if txn.AssetConfigTransaction.AssetId != 0 {
			actualId = txn.AssetConfigTransaction.AssetId
		}
		if txn.CreatedAssetIndex != 0 {
			actualId = txn.CreatedAssetIndex
		}
		err := comparisonCheck(fmt.Sprintf("asset id for returned transaction %d", idx), uint64(assetId), actualId)
		if err != nil {
			return err
		}
	}
	return nil
}

func everyTransactionIsOlderThan(olderThan string) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		olderThanDateTime, err := time.Parse(time.RFC3339, olderThan)
		if err != nil {
			return err
		}
		olderThanUnixTime := uint64(olderThanDateTime.Unix())
		if txn.RoundTime > olderThanUnixTime {
			return fmt.Errorf("txn number %d has round time of %d, but expected it to be older than unix time %d (RFC format: %s)", idx, txn.RoundTime, olderThanUnixTime, olderThan)
		}
	}
	return nil
}

func everyTransactionIsNewerThan(newerThan string) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		newerThanDateTime, err := time.Parse(time.RFC3339, newerThan)
		if err != nil {
			return err
		}
		newerThanUnixTime := uint64(newerThanDateTime.Unix())
		if txn.RoundTime < newerThanUnixTime {
			return fmt.Errorf("txn number %d has round time of %d, but expected it to be newer than unix time %d (RFC format: %s)", idx, txn.RoundTime, newerThanUnixTime, newerThan)
		}
	}
	return nil
}

func everyTransactionMovesBetweenAndCurrency(currencyGreater, currencyLesser int) error {
	for idx, txn := range indexerTransactionsResponse.Transactions {
		if txn.Type == "pay" {
			if txn.PaymentTransaction.Amount < uint64(currencyGreater) {
				return fmt.Errorf("txn number %d moved %d microAlgos, but expected it to move more than %d", idx, txn.PaymentTransaction.Amount, currencyGreater)
			}
			if currencyLesser != 0 && txn.PaymentTransaction.Amount > uint64(currencyLesser) {
				return fmt.Errorf("txn number %d moved %d microAlgos, but expected it to move less than %d", idx, txn.PaymentTransaction.Amount, currencyLesser)
			}
		}
		if txn.Type == "axfer" {
			if txn.AssetTransferTransaction.Amount < uint64(currencyGreater) {
				return fmt.Errorf("txn number %d moved %d asset units, but expected it to move more than %d", idx, txn.AssetTransferTransaction.Amount, currencyGreater)
			}
			if currencyLesser != 0 && txn.AssetTransferTransaction.Amount > uint64(currencyLesser) {
				return fmt.Errorf("txn number %d moved %d asset units, but expected it to move less than %d", idx, txn.AssetTransferTransaction.Amount, currencyLesser)
			}
		}
	}
	return nil
}

func iUseToSearchForAllTransactions(clientNum int, account string) error {
	ic := indexerClients[clientNum]
	var err error
	indexerTransactionsResponse, err = ic.LookupAccountTransactions(account).Do(context.Background())
	return err
}

func iUseToSearchForAllAssetTransactions(clientNum, assetId int) error {
	ic := indexerClients[clientNum]
	var err error
	indexerTransactionsResponse, err = ic.LookupAssetTransactions(uint64(assetId)).Do(context.Background())
	return err
}

func iGetTheNextPageUsingToSearchForTransactionsWithAnd(clientNum, limit, maxRound int) error {
	ic := indexerClients[clientNum]
	var err error
	indexerTransactionsResponse, err = ic.SearchForTransactions().Limit(uint64(limit)).MaxRound(uint64(maxRound)).NextToken(indexerTransactionsResponse.NextToken).Do(context.Background())
	return err
}

var indexerSearchForAssetsResponse []models.Asset

func iUseToSearchForAssetsWithAndToken(clientNum, zero, assetId int, creator, name, unit, token string) error {
	ic := indexerClients[clientNum]
	resp, err := ic.SearchForAssets().AssetID(uint64(assetId)).Creator(creator).Name(name).Unit(unit).NextToken(token).Do(context.Background())
	indexerSearchForAssetsResponse = resp.Assets
	return err
}

func thereAreAssetsInTheResponseTheFirstIs(numAssetsInResponse, assetIdFirstAsset int) error {
	err := comparisonCheck("number assets in search for assets response", numAssetsInResponse, len(indexerSearchForAssetsResponse))
	if err != nil {
		return err
	}
	if numAssetsInResponse == 0 {
		return nil
	}
	assetUnderScrutiny := indexerSearchForAssetsResponse[0]
	err = comparisonCheck("first asset's ID", uint64(assetIdFirstAsset), assetUnderScrutiny.Index)
	return err
}

// @indexer.applications
func iUseToSearchForApplicationsWithAndToken(indexer, limit, appId int, token string) error {
	ic := indexerClients[indexer]
	var err error
	response, err = ic.SearchForApplications().
		ApplicationId(uint64(appId)).
		Limit(uint64(limit)).
		Next(token).Do(context.Background())
	return err
}

func theParsedResponseShouldEqual(jsonfileName string) error {
	var err error
	var responseJson string

	baselinePath := path.Join("./features/resources/", jsonfileName)
	responseJson = string(json.Encode(response))

	jsonfile, err := os.Open(baselinePath)
	if err != nil {
		return err
	}
	fileBytes, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		return err
	}
	_, err = EqualJson(string(fileBytes), responseJson)
	return err
}

func iUseToLookupApplicationWith(indexer, appid int) error {
	ic := indexerClients[indexer]
	var err error
	response, err = ic.LookupApplicationByID(uint64(appid)).Do(context.Background())
	return err
}

func iUseToSearchForTransactionsWithAppIdAndToken(
	indexer, limit int, notePrefix, txType, sigType, txId string, round,
	minRound, maxRound, assetId int, beforeTime, afterTime string,
	currencyGt, currencyLt int, address, addressRole, excludeCloseTo string, appId int, token string) error {

	var err error
	ic := indexerClients[indexer]
	response, err = ic.SearchForTransactions().
		ApplicationId(uint64(appId)).
		Limit(uint64(limit)).
		Do(context.Background())
	return err
}

func iUseToSearchForAnAccountWithAppIdAndToken(indexer, assetIndex, limit, currencyGreater, currencyLesser int, arg6 string, appId int, token string) error {
	ic := indexerClients[indexer]
	var err error
	response, err = ic.SearchAccounts().
		AssetID(uint64(assetIndex)).
		Limit(uint64(limit)).
		CurrencyGreaterThan(uint64(currencyGreater)).
		CurrencyLessThan(uint64(currencyLesser)).
		ApplicationId(uint64(appId)).
		NextToken(token).Do(context.Background())
	return err
}
