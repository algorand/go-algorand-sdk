package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

var client *algod.Client
var tx types.Transaction
var transientAccount crypto.Account
var applicationId uint64

func anAlgodVClientConnectedToPortWithToken(v int, host string, port int, token string) error {
	var err error
	portAsString := strconv.Itoa(port)
	fullHost := "http://" + host + ":" + portAsString
	client, err = algod.MakeClient(fullHost, token)
	gh = []byte("MLWBXKMRJ5W3USARAFOHPQJAF4DN6KY3ZJVPIXKODKNN5ZXSZ2DQ")

	return err
}

func iCreateANewTransientAccountAndFundItWithMicroalgos(microalgos int) error {
	var err error

	transientAccount = crypto.GenerateAccount()

	status, err := acl.Status()
	if err != nil {
		return err
	}
	llr := status.LastRound

	lv, err := acl.Versions()
	gen = lv.GenesisID
	gh = lv.GenesisHash
	if err != nil {
		return err
	}

	paramsToUse := types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(llr),
		LastRoundValid:  types.Round(llr + 10),
		FlatFee:         false,
	}
	ltxn, err := future.MakePaymentTxn(accounts[1], transientAccount.Address.String(), uint64(microalgos), note, close, paramsToUse)
	if err != nil {
		return err
	}
	lsk, err := kcl.ExportKey(handle, walletPswd, accounts[1])
	if err != nil {
		return err
	}
	ltxid, lstx, err := crypto.SignTransaction(lsk.PrivateKey, ltxn)
	if err != nil {
		return err
	}
	_, err = acl.SendRawTransaction(lstx)
	if err != nil {
		return err
	}
	err = waitForTransaction(ltxid)
	if err != nil {
		return err
	}
	return nil
}

func waitForTransaction(transactionId string) error {
	status, err := acl.Status()
	if err != nil {
		return err
	}
	stopRound := status.LastRound + 10

	for {
		lstatus, err := acl.PendingTransactionInformation(transactionId)
		if err != nil {
			return err
		}
		if lstatus.ConfirmedRound > 0 {
			break // end the waiting
		}
		status, err := acl.Status()
		if err != nil {
			return err
		}
		if status.LastRound > stopRound { // last round sent tx is valid
			return fmt.Errorf("Transaction not accepted.")
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func iBuildAnApplicationTransaction(
	operation, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, appAccounts string) error {

	var clearP []byte
	var approvalP []byte
	var err error

	var ghbytes [32]byte
	copy(ghbytes[:], gh)

	var suggestedParams types.SuggestedParams
	suggestedParams, err = client.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	
	if approvalProgram != "" {
		approvalP, err = ioutil.ReadFile("features/resources/" + approvalProgram)
		if err != nil {
			return err
		}
	}

	if clearProgram != "" {
		clearP, err = ioutil.ReadFile("features/resources/" + clearProgram)
		if err != nil {
			return err
		}
	}
	args, err := parseAppArgs(appArgs)
	if err != nil {
		return err
	}
	var accs []string
	if appAccounts != "" {
		accs = strings.Split(appAccounts, ",")
	}
	fApp, err := splitUint64(foreignApps)
	if err != nil {
		return err
	}

	gSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	lSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}
	switch operation {
	case "create":
		tx, err = future.MakeApplicationCreateTx(false, approvalP, clearP,
			gSchema, lSchema, args, accs, fApp,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "update":
		tx, err = future.MakeApplicationUpdateTx(applicationId, args, accs, fApp,
			approvalP, clearP,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "call":
		tx, err = future.MakeApplicationCallTx(applicationId, args, accs,
			fApp, types.NoOpOC, approvalP, clearP, gSchema, lSchema,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "optin":
		tx, err = future.MakeApplicationOptInTx(applicationId, args, accs, fApp,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "clear":
		tx, err = future.MakeApplicationClearStateTx(applicationId, args, accs, fApp,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "closeout":
		tx, err = future.MakeApplicationCloseOutTx(applicationId, args, accs, fApp,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "delete":
		tx, err = future.MakeApplicationDeleteTx(applicationId, args, accs, fApp,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}
	}

	return nil
}

func iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs(err string) error {
	var e error
	var lstx []byte
	//	fmt.Printf("%v\n", string(json.Encode(tx)))

	txid, lstx, e = crypto.SignTransaction(transientAccount.PrivateKey, tx)
	if e != nil {
		return e
	}
	_, e = acl.SendRawTransaction(lstx)
	if e != nil {
		if strings.Contains(e.Error(), err) {
			return nil
		}
	}
	return e
}

func iWaitForTheTransactionToBeConfirmed() error {
	err := waitForTransaction(txid)
	if err != nil {
		return err
	}
	return nil
}

func iRememberTheNewApplicationID() error {

	response, err := acl.RawRequest(fmt.Sprintf("/transactions/pending/%s", txid), nil, "GET", false /* encodeJSON */, nil)
	if err != nil {
		return err
	}
	if txres := response["txresults"]; txres != nil {
		createdapp := txres.(map[string]interface{})["createdapp"]
		applicationId = uint64(createdapp.(float64))
	}
	return nil
}

func parseAppArgs(appArgsString string) (appArgs [][]byte, err error) {
	if appArgsString == "" {
		return make([][]byte, 0), nil
	}
	argsArray := strings.Split(appArgsString, ",")
	resp := make([][]byte, len(argsArray))

	for idx, arg := range argsArray {
		typeArg := strings.Split(arg, ":")
		switch typeArg[0] {
		case "str":
			resp[idx] = []byte(typeArg[1])
		default:
			return nil, fmt.Errorf("Applications doesn't currently support argument of type %s", typeArg[0])
		}
	}
	return resp, err
}

func splitUint64(uint64s string) ([]uint64, error) {
	if uint64s == "" {
		return make([]uint64, 0), nil
	}
	uintarr := strings.Split(uint64s, ",")
	resp := make([]uint64, len(uintarr))
	var err error
	for idx, val := range uintarr {
		resp[idx], err = strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func ApplicationsContext(s *godog.Suite) {
	s.Step(`^an algod v(\d+) client connected to "([^"]*)" port (\d+) with token "([^"]*)"$`, anAlgodVClientConnectedToPortWithToken)
	s.Step(`^I create a new transient account and fund it with (\d+) microalgos\.$`, iCreateANewTransientAccountAndFundItWithMicroalgos)
	s.Step(`^I build an application transaction with the transient account, the current application, suggested params, operation "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", app-accounts "([^"]*)"$`, iBuildAnApplicationTransaction)
	s.Step(`^I sign and submit the transaction, saving the txid\. If there is an error it is "([^"]*)"\.$`, iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs)
	s.Step(`^I wait for the transaction to be confirmed\.$`, iWaitForTheTransactionToBeConfirmed)
	s.Step(`^I remember the new application ID\.$`, iRememberTheNewApplicationID)
}
