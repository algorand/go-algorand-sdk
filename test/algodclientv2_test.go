package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/cucumber/godog"
)

func AlgodClientV2Context(s *godog.Suite) {
	s.Step(`^mock http responses in "([^"]*)" loaded from "([^"]*)"$`, mockHttpResponsesInLoadedFrom)
	s.Step(`^expect error string to contain "([^"]*)"$`, expectErrorStringToContain)
	s.Step(`^we make any Pending Transaction Information call$`, weMakeAnyPendingTransactionInformationCall)
	s.Step(`^the parsed Pending Transaction Information response should have sender "([^"]*)"$`, theParsedPendingTransactionInformationResponseShouldHaveSender)
	s.Step(`^we make any Pending Transactions Information call$`, weMakeAnyPendingTransactionsInformationCall)
	s.Step(`^the parsed Pending Transactions Information response should have sender "([^"]*)"$`, theParsedPendingTransactionsInformationResponseShouldHaveSender)
	s.Step(`^we make any Send Raw Transaction call$`, weMakeAnySendRawTransactionCall)
	s.Step(`^the parsed Send Raw Transaction response should have txid "([^"]*)"$`, theParsedSendRawTransactionResponseShouldHaveTxid)
	s.Step(`^we make any Pending Transactions By Address call$`, weMakeAnyPendingTransactionsByAddressCall)
	s.Step(`^the parsed Pending Transactions By Address response should contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedPendingTransactionsByAddressResponseShouldContainAnArrayOfLenAndElementNumberShouldHaveSender)
	s.Step(`^we make any Node Status call$`, weMakeAnyNodeStatusCall)
	s.Step(`^the parsed Node Status response should have a last round of (\d+)$`, theParsedNodeStatusResponseShouldHaveALastRoundOf)
	s.Step(`^we make any Ledger Supply call$`, weMakeAnyLedgerSupplyCall)
	s.Step(`^the parsed Ledger Supply response should have totalMoney (\d+) onlineMoney (\d+) on round (\d+)$`, theParsedLedgerSupplyResponseShouldHaveTotalMoneyOnlineMoneyOnRound)
	s.Step(`^we make any Status After Block call$`, weMakeAnyStatusAfterBlockCall)
	s.Step(`^the parsed Status After Block response should have a last round of (\d+)$`, theParsedStatusAfterBlockResponseShouldHaveALastRoundOf)
	s.Step(`^we make any Account Information call$`, weMakeAnyAccountInformationCall)
	s.Step(`^the parsed Account Information response should have address "([^"]*)"$`, theParsedAccountInformationResponseShouldHaveAddress)
	s.Step(`^we make any Get Block call$`, weMakeAnyGetBlockCall)
	s.Step(`^the parsed Get Block response should have rewards pool "([^"]*)"$`, theParsedGetBlockResponseShouldHaveRewardsPool)
	s.Step(`^we make any Suggested Transaction Parameters call$`, weMakeAnySuggestedTransactionParametersCall)
	s.Step(`^the parsed Suggested Transaction Parameters response should have first round valid of (\d+)$`, theParsedSuggestedTransactionParametersResponseShouldHaveFirstRoundValidOf)
	s.Step(`^expect the path used to be "([^"]*)"$`, expectThePathUsedToBe)
	s.Step(`^we make a Pending Transaction Information against txid "([^"]*)" with max (\d+)$`, weMakeAPendingTransactionInformationAgainstTxidWithMax)
	s.Step(`^we make a Pending Transactions By Address call against account "([^"]*)" and max (\d+)$`, weMakeAPendingTransactionsByAddressCallAgainstAccountAndMax)
	s.Step(`^we make a Status after Block call with round (\d+)$`, weMakeAStatusAfterBlockCallWithRound)
	s.Step(`^we make an Account Information call against account "([^"]*)"$`, weMakeAnAccountInformationCallAgainstAccount)
	s.Step(`^we make a Get Block call against block number (\d+)$`, weMakeAGetBlockCallAgainstBlockNumber)
	s.Step(`^the parsed Pending Transactions Information response should contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedPendingTransactionsInformationResponseShouldContainAnArrayOfLenAndElementNumberShouldHaveSender)
	s.Step(`^we make a Pending Transaction Information against txid "([^"]*)" with format "([^"]*)"$`, weMakeAPendingTransactionInformationAgainstTxidWithFormat)
	s.Step(`^we make a Pending Transaction Information with max (\d+) and format "([^"]*)"$`, weMakeAPendingTransactionInformationWithMaxAndFormat)
	s.Step(`^we make a Pending Transactions By Address call against account "([^"]*)" and max (\d+) and format "([^"]*)"$`, weMakeAPendingTransactionsByAddressCallAgainstAccountAndMaxAndFormat)
	s.Step(`^we make a Get Block call against block number (\d+) with format "([^"]*)"$`, weMakeAGetBlockCallAgainstBlockNumberWithFormat)
	s.BeforeScenario(func(interface{}) {
		globalErrForExamination = nil
	})
}

var stxResponse types.SignedTxn
var pendingTransactionInformationResponse models.PendingTransactionInfoResponse

func weMakeAnyPendingTransactionInformationCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	pendingTransactionInformationResponse, stxResponse, globalErrForExamination = algodClient.PendingTransactionInformation("").Do(context.Background())
	return nil
}

func theParsedPendingTransactionInformationResponseShouldHaveSender(sender string) error {
	if stxResponse.Txn.Sender.String() != sender {
		return fmt.Errorf("expected txn to have sender %s but actual sender was %s", sender, stxResponse.Txn.Sender.String())
	}
	return nil
}

var stxsResponse []types.SignedTxn

func weMakeAnyPendingTransactionsInformationCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, stxsResponse, globalErrForExamination = algodClient.PendingTransactions().Do(context.Background())
	return nil
}

func theParsedPendingTransactionsInformationResponseShouldHaveSender(sender string) error {
	if stxsResponse[0].Txn.Sender.String() != sender {
		return fmt.Errorf("expected txn to have sender %s but actual sender was %s", sender, stxsResponse[0].Txn.Sender.String())
	}
	return nil
}

func theParsedPendingTransactionsInformationResponseShouldContainAnArrayOfLenAndElementNumberShouldHaveSender(length, idx int, sender string) error {
	if len(stxsResponse) != length {
		return fmt.Errorf("expected response length %d but received length %d", length, len(stxsResponse))
	}
	if stxsResponse[idx].Txn.Sender.String() != sender {
		return fmt.Errorf("expected txn %d to have sender %s but real sender was %s", idx, sender, stxsResponse[idx].Txn.Sender.String())
	}
	return nil
}

var txidResponse string

func weMakeAnySendRawTransactionCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	txidResponse, globalErrForExamination = algodClient.SendRawTransaction(nil).Do(context.Background())
	return nil
}

func theParsedSendRawTransactionResponseShouldHaveTxid(txid string) error {
	if txidResponse != txid {
		return fmt.Errorf("expected txn to have txid %s but actual txid was %s", txidResponse, txid)
	}
	return nil
}

func weMakeAnyPendingTransactionsByAddressCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, stxsResponse, globalErrForExamination = algodClient.PendingTransactionsByAddress("").Do(context.Background())
	return nil
}

func theParsedPendingTransactionsByAddressResponseShouldContainAnArrayOfLenAndElementNumberShouldHaveSender(expectedLen, idx int, expectedSender string) error {
	length := len(stxsResponse)
	if length != expectedLen {
		return fmt.Errorf("length of response %d mismatched expected length %d", length, expectedLen)
	}
	if stxsResponse[idx].Txn.Sender.String() != expectedSender {
		return fmt.Errorf("response sender %s mismatched expected sender %s", stxsResponse[idx].Txn.Sender.String(), expectedSender)
	}
	return nil
}

var statusResponse models.NodeStatus

func weMakeAnyNodeStatusCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	statusResponse, globalErrForExamination = algodClient.Status().Do(context.Background())
	return nil
}

func theParsedNodeStatusResponseShouldHaveALastRoundOf(lastRound int) error {
	if statusResponse.LastRound != uint64(lastRound) {
		return fmt.Errorf("response last round %d mismatched expected last round %d", statusResponse.LastRound, lastRound)
	}
	return nil
}

var supplyResponse models.Supply

func weMakeAnyLedgerSupplyCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	supplyResponse, globalErrForExamination = algodClient.Supply().Do(context.Background())
	return nil
}

func theParsedLedgerSupplyResponseShouldHaveTotalMoneyOnlineMoneyOnRound(total, online, round int) error {
	if supplyResponse.TotalMoney != uint64(total) {
		return fmt.Errorf("response total money %d mismatched expected total %d", supplyResponse.TotalMoney, uint64(total))
	}
	if supplyResponse.OnlineMoney != uint64(online) {
		return fmt.Errorf("response online money %d mismatched expected online money %d", supplyResponse.OnlineMoney, uint64(online))
	}
	if supplyResponse.Round != uint64(round) {
		return fmt.Errorf("response round %d mismatched expected round %d", supplyResponse.Round, uint64(round))
	}
	return nil
}

func weMakeAnyStatusAfterBlockCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	statusResponse, globalErrForExamination = algodClient.StatusAfterBlock(0).Do(context.Background())
	return nil
}

func theParsedStatusAfterBlockResponseShouldHaveALastRoundOf(lastRound int) error {
	if statusResponse.LastRound != uint64(lastRound) {
		return fmt.Errorf("response last round %d mismatched expected last round %d", statusResponse.LastRound, lastRound)
	}
	return nil
}

var accountResponse models.Account

func weMakeAnyAccountInformationCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	accountResponse, globalErrForExamination = algodClient.AccountInformation("").Do(context.Background())
	return nil
}

func theParsedAccountInformationResponseShouldHaveAddress(address string) error {
	if accountResponse.Address != address {
		return fmt.Errorf("response address %s mismatched expected address %s", accountResponse.Address, address)
	}
	return nil
}

var blockResponse types.Block

func weMakeAnyGetBlockCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	blockResponse, globalErrForExamination = algodClient.Block(0).Do(context.Background())
	return nil
}

func theParsedGetBlockResponseShouldHaveRewardsPool(pool string) error {
	blockResponseRewardsPoolBytes := [32]byte(blockResponse.RewardsPool)
	poolBytes, err := base64.StdEncoding.DecodeString(pool)
	if err != nil {
		return err
	}
	if !bytes.Equal(poolBytes, blockResponseRewardsPoolBytes[:]) {
		return fmt.Errorf("response pool %v mismatched expected pool %v", blockResponseRewardsPoolBytes, poolBytes)
	}
	return nil
}

var suggestedParamsResponse types.SuggestedParams

func weMakeAnySuggestedTransactionParametersCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	suggestedParamsResponse, globalErrForExamination = algodClient.SuggestedParams().Do(context.Background())
	return nil
}

func theParsedSuggestedTransactionParametersResponseShouldHaveFirstRoundValidOf(firstValid int) error {
	if suggestedParamsResponse.FirstRoundValid != types.Round(firstValid) {
		return fmt.Errorf("response first round valid %d mismatched expected first round valid %d", suggestedParamsResponse.FirstRoundValid, types.Round(firstValid))
	}
	return nil
}

func weMakeAPendingTransactionInformationAgainstTxidWithMax(txid string, max int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactionInformation(txid).Max(uint64(max)).Do(context.Background())
	return nil
}

func weMakeAPendingTransactionsByAddressCallAgainstAccountAndMax(account string, max int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactionsByAddress(account).Max(uint64(max)).Do(context.Background())
	return nil
}

func weMakeAStatusAfterBlockCallWithRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = algodClient.StatusAfterBlock(uint64(round)).Do(context.Background())
	return nil
}

func weMakeAnAccountInformationCallAgainstAccount(account string) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = algodClient.AccountInformation(account).Do(context.Background())
	return nil
}

func weMakeAGetBlockCallAgainstBlockNumber(blocknum int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = algodClient.Block(uint64(blocknum)).Do(context.Background())
	return nil
}

func weMakeAPendingTransactionInformationAgainstTxidWithFormat(txid, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	return weMakeAPendingTransactionInformationAgainstTxidWithMax(txid, 0)
}

func weMakeAPendingTransactionInformationWithMaxAndFormat(max int, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactions().Max(uint64(max)).Do(context.Background())
	return nil
}

func weMakeAPendingTransactionsByAddressCallAgainstAccountAndMaxAndFormat(account string, max int, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	return weMakeAPendingTransactionsByAddressCallAgainstAccountAndMax(account, max)
}

func weMakeAGetBlockCallAgainstBlockNumberWithFormat(blocknum int, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	return weMakeAGetBlockCallAgainstBlockNumber(blocknum)
}
