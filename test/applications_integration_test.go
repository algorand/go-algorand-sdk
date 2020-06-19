package main

import (
	"fmt"
	"strconv"
	"io/ioutil"
	
	"github.com/cucumber/godog"
	
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

var client algod.Client
var tx types.Transaction 

func anAlgodVClientConnectedToPortWithToken(v int, host string, port int, token string) error {
	portAsString := strconv.Itoa(port)
	fullHost := "http://" + host + ":" + portAsString
	client, err := algod.MakeClient(fullHost, token)
	return err
}

func iCreateANewTransientAccountAndFundItWithMicroalgos(microalgos int) error {
	return godog.ErrPending
}

func iBuildAnApplicationTransactionWithTheTransientAccountTheCurrentApplicationSuggestedParamsOperationApprovalprogramClearprogramGlobalbytesGlobalintsLocalbytesLocalintsAppargsForeignappsAppaccounts(
	operation, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, appAccounts string) error {

	approvalP, err := ioutil.ReadFile(approvalProgram)
	clearP, err := ioutil.ReadFile(clearProgram)
	
	switch (operation) {
	case "create":
		tx, err := transaction.MakeUnsignedAppCreateTx(0, approvalP, clearP,
			types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}, 
			types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)},
			appArgs, appAccounts, foreignApps)
	case "update":

	case "call":

	case "optin":

	case "clear":

	case "closeout":

	case "delete":

		
	}
	
	return nil
}

func iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs(err string) error {
	return godog.ErrPending
}

func iWaitForTheTransactionToBeConfirmed() error {
	return godog.ErrPending
}

func iRememberTheNewApplicationID() error {
	return godog.ErrPending
}

func ApplicationsContext(s *godog.Suite) {
	s.Step(`^an algod v(\d+) client connected to "([^"]*)" port (\d+) with token "([^"]*)"$`, anAlgodVClientConnectedToPortWithToken)
	s.Step(`^I create a new transient account and fund it with (\d+) microalgos\.$`, iCreateANewTransientAccountAndFundItWithMicroalgos)
	s.Step(`^I build an application transaction with the transient account, the current application, suggested params, operation "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", app-accounts "([^"]*)"$`, iBuildAnApplicationTransactionWithTheTransientAccountTheCurrentApplicationSuggestedParamsOperationApprovalprogramClearprogramGlobalbytesGlobalintsLocalbytesLocalintsAppargsForeignappsAppaccounts)
	s.Step(`^I sign and submit the transaction, saving the txid\. If there is an error it is "([^"]*)"\.$`, iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs)
	s.Step(`^I wait for the transaction to be confirmed\.$`, iWaitForTheTransactionToBeConfirmed)
	s.Step(`^I remember the new application ID\.$`, iRememberTheNewApplicationID)
}

func main () {
	fmt.Println("bla")
}
