package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/algorand/go-algorand-sdk/v2client/algod"
	"github.com/algorand/go-algorand-sdk/v2client/common/models"
)

func AlgodClientV2Context(s *godog.Suite) {
	s.Step(`^we make any Shutdown call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyShutdownCallReturnMockResponse)
	s.Step(`^we make any Register Participation Keys call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyRegisterParticipationKeysCallReturnMockResponse)
	s.Step(`^we make any Pending Transaction Information call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyPendingTransactionInformationCallReturnMockResponse)
	s.Step(`^the parsed Pending Transaction Information response should have sender "([^"]*)"$`, theParsedPendingTransactionInformationResponseShouldHaveSender)
	s.Step(`^we make any Send Raw Transaction call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnySendRawTransactionCallReturnMockResponse)
	s.Step(`^the parsed Send Raw Transaction response should have txid "([^"]*)"$`, theParsedSendRawTransactionResponseShouldHaveTxid)
	s.Step(`^we make any Pending Transactions By Address call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyPendingTransactionsByAddressCallReturnMockResponse)
	s.Step(`^the parsed Pending Transactions By Address response should contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedPendingTransactionsByAddressResponseShouldContainAnArrayOfLenAndElementNumberShouldHaveSender)
	s.Step(`^we make any Node Status call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyNodeStatusCallReturnMockResponse)
	s.Step(`^the parsed Node Status response should have a last round of (\d+)$`, theParsedNodeStatusResponseShouldHaveALastRoundOf)
	s.Step(`^we make any Ledger Supply call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyLedgerSupplyCallReturnMockResponse)
	s.Step(`^the parsed Ledger Supply response should have totalMoney (\d+) onlineMoney (\d+) on round (\d+)$`, theParsedLedgerSupplyResponseShouldHaveTotalMoneyOnlineMoneyOnRound)
	s.Step(`^we make any Status After Block call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyStatusAfterBlockCallReturnMockResponse)
	s.Step(`^the parsed Status After Block response should have a last round of (\d+)$`, theParsedStatusAfterBlockResponseShouldHaveALastRoundOf)
	s.Step(`^we make any Account Information call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyAccountInformationCallReturnMockResponse)
	s.Step(`^the parsed Account Information response should have address "([^"]*)"$`, theParsedAccountInformationResponseShouldHaveAddress)
	s.Step(`^we make any Get Block call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnyGetBlockCallReturnMockResponse)
	s.Step(`^the parsed Get Block response should have proposer "([^"]*)"$`, theParsedGetBlockResponseShouldHaveProposer)
	s.Step(`^we make any Suggested Transaction Parameters call, return mock response "([^"]*)" and expect error string to contain "([^"]*)"$`, weMakeAnySuggestedTransactionParametersCallReturnMockResponse)
	s.Step(`^the parsed Suggested Transaction Parameters response should have first round valid of (\d+)$`, theParsedSuggestedTransactionParametersResponseShouldHaveFirstRoundValidOf)
	s.BeforeScenario(func(interface{}) {
	})
}

func buildMockAlgodv2AndClient(jsonfile string) (*httptest.Server, algod.Client, error) {
	jsonBytes, err := loadMockJson(jsonfile)
	if err != nil {
		return nil, algod.Client{}, err
	}
	mockAlgod := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}))
	noToken := ""
	algodClient, err := algod.MakeClient(mockAlgod.URL, noToken)
	return mockAlgod, algodClient, err
}

func weMakeAnyShutdownCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	err = algodClient.Shutdown(context.Background(), models.ShutdownParams{})
	return confirmErrorContainsString(err, errorString)
}

func weMakeAnyRegisterParticipationKeysCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	err = algodClient.RegisterParticipationKeys(context.Background(), "", models.RegisterParticipationKeysAccountIdParams{})
	return err
}

var transactionResult types.Transaction

func weMakeAnyPendingTransactionInformationCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	transactionResult, err = algodClient.PendingTransactionInformation(context.Background(), "", models.GetPendingTransactionsParams{})
	return confirmErrorContainsString(err, errorString)
}

func theParsedPendingTransactionInformationResponseShouldHaveSender(sender string) error {
	if transactionResult.Sender.String() != sender {
		return fmt.Errorf("decoded sender %s mismatched with expected sender %s", transactionResult.Sender.String(), sender)
	}
	return nil
}

var txIdResult string

func weMakeAnySendRawTransactionCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	txIdResult, err = algodClient.SendRawTransaction(context.Background(), []byte{})
	return confirmErrorContainsString(err, errorString)
}

func theParsedSendRawTransactionResponseShouldHaveTxid(txid string) error {
	if txIdResult != txid {
		return fmt.Errorf("decoded txid %s mismatched with expected txid %s", txIdResult, txid)
	}
	return nil
}

var transactionResults []types.Transaction

func weMakeAnyPendingTransactionsByAddressCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	transactionResults, err = algodClient.PendingTransactionsByAddress(context.Background(), "", models.GetPendingTransactionsByAddressParams{})
	return confirmErrorContainsString(err, errorString)
}

func theParsedPendingTransactionsByAddressResponseShouldContainAnArrayOfLenAndElementNumberShouldHaveSender(length, idx int, sender string) error {
	if len(transactionResults) != length {
		return fmt.Errorf("decoded response length %d mismatched with expected array length %d", len(transactionResults), length)
	}
	if transactionResults[idx].Sender.String() != sender {
		return fmt.Errorf("decoded sender %s mismatched with expected sender %s", transactionResult.Sender.String(), sender)
	}
	return nil
}

var statusResult models.NodeStatus

func weMakeAnyNodeStatusCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	statusResult, err = algodClient.Status(context.Background())
	return confirmErrorContainsString(err, errorString)
}

func theParsedNodeStatusResponseShouldHaveALastRoundOf(roundNum int) error {
	if statusResult.LastRound != uint64(roundNum) {
		return fmt.Errorf("decoded status last round %d mismatched with expected round number %d", statusResult.LastRound, roundNum)
	}
	return nil
}

var supplyResult models.Supply

func weMakeAnyLedgerSupplyCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	supplyResult, err = algodClient.Supply(context.Background())
	return confirmErrorContainsString(err, errorString)
}

func theParsedLedgerSupplyResponseShouldHaveTotalMoneyOnlineMoneyOnRound(totalMoney, onlineMoney, roundNum int) error {
	if supplyResult.TotalMoney != uint64(totalMoney) {
		return fmt.Errorf("decoded supply total money %d mismatched with expected total money %d", supplyResult.TotalMoney, totalMoney)
	}
	if supplyResult.OnlineMoney != uint64(onlineMoney) {
		return fmt.Errorf("decoded supply online money %d mismatched with expected online money %d", supplyResult.OnlineMoney, onlineMoney)
	}
	if supplyResult.Round != uint64(roundNum) {
		return fmt.Errorf("decoded supply round %d mismatched with expected round number %d", supplyResult.Round, roundNum)
	}
	return nil
}

func weMakeAnyStatusAfterBlockCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	statusResult, err = algodClient.StatusAfterBlock(context.Background(), 1)
	return confirmErrorContainsString(err, errorString)
}

func theParsedStatusAfterBlockResponseShouldHaveALastRoundOf(roundNum int) error {
	if statusResult.LastRound != uint64(roundNum) {
		return fmt.Errorf("decoded status last round %d mismatched with expected round number %d", statusResult.LastRound, roundNum)
	}
	return nil
}

var accountResult models.Account

func weMakeAnyAccountInformationCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	accountResult, err = algodClient.AccountInformation(context.Background(), "")
	return confirmErrorContainsString(err, errorString)
}

func theParsedAccountInformationResponseShouldHaveAddress(address string) error {
	if accountResult.Address != address {
		return fmt.Errorf("parsed address %s mismatched with expected address %s", accountResult.Address, address)
	}
	return nil
}

var blockResult models.Block

func weMakeAnyGetBlockCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	blockResult, err = algodClient.Block(context.Background(), 0, models.GetBlockParams{})
	return confirmErrorContainsString(err, errorString)
}

func theParsedGetBlockResponseShouldHaveProposer(proposer string) error {
	if blockResult.Proposer != proposer {
		return fmt.Errorf("parsed proposer %s mismatched with expected address %s", blockResult.Proposer, proposer)
	}
	return nil
}

var suggestedParamsResult types.SuggestedParams

func weMakeAnySuggestedTransactionParametersCallReturnMockResponse(jsonfile string, errorString string) error {
	mockAlgod, algodClient, err := buildMockAlgodv2AndClient(jsonfile)
	if mockAlgod != nil {
		defer mockAlgod.Close()
	}
	if err != nil {
		return err
	}
	suggestedParamsResult, err = algodClient.SuggestedParams(context.Background())
	return confirmErrorContainsString(err, errorString)
}

func theParsedSuggestedTransactionParametersResponseShouldHaveFirstRoundValidOf(firstRoundValid int) error {
	if suggestedParamsResult.FirstRoundValid != types.Round(firstRoundValid) {
		return fmt.Errorf("decoded suggested params first round valid %d mismatched with expected round number %d", suggestedParamsResult.FirstRoundValid, firstRoundValid)
	}
	return nil
}
