package test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cucumber/godog"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

func weMakeASearchForApplicationsCallWithAnd(arg1, arg2 int) error {
	return godog.ErrPending
}

func weExpectThePathUsedToBe(arg1 string) error {
	return godog.ErrPending
}

func weMakeALookupApplicationsCallWithAnd(arg1, arg2 int) error {
	return godog.ErrPending
}

var sk1 ed25519.PrivateKey
var addr1 types.Address

func aSigningAccountWithAddressAndMnemonic(address, mnem string) error {
	var err error
	addr1, err = types.DecodeAddress(address)
	if err != nil {
		return err
	}
	sk1, err = mnemonic.ToPrivateKey(mnem)
	return err
}

func iBuildAnApplicationTransactionUnit(
	operation string,
	applicationIdInt int,
	sender, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, appAccounts string,
	fee, firstValid, lastValid int,
	genesisHash string) error {

	applicationId = uint64(applicationIdInt)
	var clearP []byte
	var approvalP []byte
	var err error

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
		tx, err = transaction.MakeUnsignedAppCreateTx(types.NoOpOC, approvalP, clearP,
			gSchema, lSchema, args, accs, fApp)
		if err != nil {
			return err
		}

	case "update":
		tx, err = transaction.MakeUnsignedAppUpdateTx(applicationId, args, accs, fApp,
			approvalP, clearP)
		if err != nil {
			return err
		}

	case "call":
		tx, err = transaction.MakeUnsignedApplicationCallTx(applicationId, args, accs,
			fApp, types.NoOpOC, approvalP, clearP, gSchema, lSchema)
	case "optin":
		tx, err = transaction.MakeUnsignedAppOptInTx(applicationId, args, accs, fApp)
		if err != nil {
			return err
		}

	case "clear":
		tx, err = transaction.MakeUnsignedAppClearStateTx(applicationId, args, accs, fApp)
		if err != nil {
			return err
		}

	case "closeout":
		tx, err = transaction.MakeUnsignedAppCloseOutTx(applicationId, args, accs, fApp)
		if err != nil {
			return err
		}

	case "delete":
		tx, err = transaction.MakeUnsignedAppDeleteTx(applicationId, args, accs, fApp)
		if err != nil {
			return err
		}
	}

	suggestedParams, err := getSuggestedParams(uint64(fee), uint64(firstValid), uint64(lastValid), "", genesisHash, true)
	if err != nil {
		return err
	}

	transaction.SetApplicationTransactionFields(
		&tx, //tx types.Transaction,
		suggestedParams,
		addr1,           //sender types.Address,
		nil,             //note []byte,
		types.Digest{},  //group types.Digest,
		[32]byte{},      //lease [32]byte,
		types.Address{}) //rekeyTo types.Address
	return nil

}

func signTheTransaction() error {
	var err error
	_, stx, err = crypto.SignTransaction(sk1, tx)
	return err
}

func theBaseEncodedSignedTransactionShouldEqual(base int, golden string) error {
	gold, err := base64.StdEncoding.DecodeString(golden)
	if err != nil {
		return err
	}
	if !bytes.Equal(gold, stx) {
		return fmt.Errorf("Application signed transaction does not match the golden.")
	}
	return nil
}

func getSuggestedParams(
	fee, fv, lv uint64,
	gen, ghb64 string,
	flat bool) (types.SuggestedParams, error) {
	gh, err := base64.StdEncoding.DecodeString(ghb64)
	if err != nil {
		return types.SuggestedParams{}, err
	}
	return types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         flat,
	}, err
}

func ApplicationsUnitContext(s *godog.Suite) {
	s.Step(`^a signing account with address "([^"]*)" and mnemonic "([^"]*)"$`, aSigningAccountWithAddressAndMnemonic)
	s.Step(`^I build an application transaction with operation "([^"]*)", application-id (\d+), sender "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", app-accounts "([^"]*)", fee (\d+), first-valid (\d+), last-valid (\d+), genesis-hash "([^"]*)"$`, iBuildAnApplicationTransactionUnit)
	s.Step(`^sign the transaction$`, signTheTransaction)
	s.Step(`^the base(\d+) encoded signed transaction should equal "([^"]*)"$`, theBaseEncodedSignedTransactionShouldEqual)
}
