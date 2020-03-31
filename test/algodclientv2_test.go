package main

import (
	"github.com/cucumber/godog"
)

func AlgodClientV2Context(s *godog.Suite) {
	s.Step(`^mock http responses in "([^"]*)" loaded from "([^"]*)"$`, mockHttpResponsesInLoadedFrom)
	s.Step(`^we make any Shutdown call$`, weMakeAnyShutdownCall)
	s.Step(`^expect error string to contain "([^"]*)"$`, expectErrorStringToContain)
	s.Step(`^we make any Register Participation Keys call$`, weMakeAnyRegisterParticipationKeysCall)
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
	s.Step(`^the parsed Get Block response should have proposer "([^"]*)"$`, theParsedGetBlockResponseShouldHaveProposer)
	s.Step(`^we make any Suggested Transaction Parameters call$`, weMakeAnySuggestedTransactionParametersCall)
	s.Step(`^the parsed Suggested Transaction Parameters response should have first round valid of (\d+)$`, theParsedSuggestedTransactionParametersResponseShouldHaveFirstRoundValidOf)
	s.BeforeScenario(func(interface{}) {})
}

//
//func buildMockAlgodv2AndClient(jsonfile string) (*httptest.Server, algod.Client, error) {
//	jsonBytes, err := loadMockJsons("", jsonfile)
//	if err != nil {
//		return nil, algod.Client{}, err
//	}
//	mockAlgod := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		w.Write(jsonBytes[0])
//	}))
//	noToken := ""
//	algodClient, err := algod.MakeClient(mockAlgod.URL, noToken)
//	return mockAlgod, algodClient, err
//}

func weMakeAnyShutdownCall() error {
	return godog.ErrPending
}

func expectErrorStringToContain(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnyRegisterParticipationKeysCall() error {
	return godog.ErrPending
}

func weMakeAnyPendingTransactionInformationCall() error {
	return godog.ErrPending
}

func theParsedPendingTransactionInformationResponseShouldHaveSender(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnyPendingTransactionsInformationCall() error {
	return godog.ErrPending
}

func theParsedPendingTransactionsInformationResponseShouldHaveSender(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnySendRawTransactionCall() error {
	return godog.ErrPending
}

func theParsedSendRawTransactionResponseShouldHaveTxid(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnyPendingTransactionsByAddressCall() error {
	return godog.ErrPending
}

func theParsedPendingTransactionsByAddressResponseShouldContainAnArrayOfLenAndElementNumberShouldHaveSender(arg1, arg2 int, arg3 string) error {
	return godog.ErrPending
}

func weMakeAnyNodeStatusCall() error {
	return godog.ErrPending
}

func theParsedNodeStatusResponseShouldHaveALastRoundOf(arg1 int) error {
	return godog.ErrPending
}

func weMakeAnyLedgerSupplyCall() error {
	return godog.ErrPending
}

func theParsedLedgerSupplyResponseShouldHaveTotalMoneyOnlineMoneyOnRound(arg1, arg2, arg3 int) error {
	return godog.ErrPending
}

func weMakeAnyStatusAfterBlockCall() error {
	return godog.ErrPending
}

func theParsedStatusAfterBlockResponseShouldHaveALastRoundOf(arg1 int) error {
	return godog.ErrPending
}

func weMakeAnyAccountInformationCall() error {
	return godog.ErrPending
}

func theParsedAccountInformationResponseShouldHaveAddress(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnyGetBlockCall() error {
	return godog.ErrPending
}

func theParsedGetBlockResponseShouldHaveProposer(arg1 string) error {
	return godog.ErrPending
}

func weMakeAnySuggestedTransactionParametersCall() error {
	return godog.ErrPending
}

func theParsedSuggestedTransactionParametersResponseShouldHaveFirstRoundValidOf(arg1 int) error {
	return godog.ErrPending
}
