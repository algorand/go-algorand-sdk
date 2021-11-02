package test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/cucumber/godog"
)

var signingAccount crypto.Account
var txnSuggestedParams types.SuggestedParams

func aSigningAccountWithAddressAndMnemonic(address, mnem string) error {
	sk, err := mnemonic.ToPrivateKey(mnem)
	if err != nil {
		return err
	}

	signingAccount, err = crypto.AccountFromPrivateKey(sk)
	if err != nil {
		return err
	}

	addr, err := types.DecodeAddress(address)
	if err != nil {
		return err
	}

	if addr != signingAccount.Address {
		return fmt.Errorf("Supplied address does not match private key: %s != %s", addr.String(), signingAccount.Address.String())
	}

	return err
}

func suggestedTransactionParameters(fee int, flatFee string, firstValid, lastValid int, genesisHash, genesisId string) error {
	if fee < 0 {
		return fmt.Errorf("Fee is negative: %d", fee)
	}
	if firstValid < 0 {
		return fmt.Errorf("First valid is negative: %d", firstValid)
	}
	if lastValid < 0 {
		return fmt.Errorf("Last valid is negative: %d", lastValid)
	}
	flatFeeBool, err := strconv.ParseBool(flatFee)
	if err != nil {
		return fmt.Errorf("Could not parse flatFee boolean value: %v", err)
	}

	genesisHashBytes, err := base64.StdEncoding.DecodeString(genesisHash)
	if err != nil {
		return err
	}

	txnSuggestedParams = types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		FlatFee:         flatFeeBool,
		FirstRoundValid: types.Round(firstValid),
		LastRoundValid:  types.Round(lastValid),
		GenesisHash:     genesisHashBytes,
		GenesisID:       genesisId,
	}

	return nil
}

func signTheTransaction() error {
	var err error
	_, stx, err = crypto.SignTransaction(signingAccount.PrivateKey, tx)
	return err
}

func theBase64EncodedSignedTransactionShouldEqual(golden string) error {
	gold, err := base64.StdEncoding.DecodeString(golden)
	if err != nil {
		return err
	}
	if !bytes.Equal(gold, stx) {
		toPrint := base64.StdEncoding.EncodeToString(stx)
		return fmt.Errorf("Actual signed transaction does not match the expected. Got %s", toPrint)
	}
	return nil
}

func theDecodedTransactionShouldEqualTheOriginal() error {
	var signedTxn types.SignedTxn
	err := msgpack.Decode(stx, &signedTxn)
	if err != nil {
		return fmt.Errorf("Could not decode signed transaction: %v", err)
	}

	if !reflect.DeepEqual(signedTxn.Txn, tx) {
		return fmt.Errorf("Decoded signed transaction does not equal the original")
	}

	return nil
}

func buildKeyregTransaction(sender, nonparticipation string,
	voteFirst, voteLast, keyDilution int,
	votePkB64, selectionPkB64, stateProofPkB64 string) error {

	if voteFirst < 0 || voteLast < 0 || keyDilution < 0 {
		return fmt.Errorf("Integer arguments cannot be negative")
	}

	nonpart, err := strconv.ParseBool(nonparticipation)
	if err != nil {
		return fmt.Errorf("Could not parse nonparticipation value: %v", err)
	}

	var stateProofPk types.Verifier
	if len(stateProofPkB64) > 0 {
		stateProofBytes, err := base64.StdEncoding.DecodeString(stateProofPkB64)
		if err != nil {
			return err
		}
		if len(stateProofBytes) != len(stateProofPk.Root) {
			return fmt.Errorf("Invalid length for state proof key. Expected %d, got %d", len(stateProofPk.Root), len(stateProofBytes))
		}
		copy(stateProofPk.Root[:], stateProofBytes)
	}

	tx, err = future.MakeKeyRegTxnV2(sender, nil, txnSuggestedParams, votePkB64, selectionPkB64, uint64(voteFirst), uint64(voteLast), uint64(keyDilution), nonpart, stateProofPk)
	return err
}

// this function is a legacy step, we should get rid of it in favor one a newer step that works like
// buildKeyregTransaction above
func buildLegacyAppCallTransaction(
	operation string,
	applicationIdInt int,
	sender, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, foreignAssets, appAccounts string,
	fee, firstValid, lastValid int,
	genesisHash string, extraPages int) error {

	if applicationIdInt < 0 || globalBytes < 0 || globalInts < 0 || localBytes < 0 || localInts < 0 || extraPages < 0 || fee < 0 || firstValid < 0 || lastValid < 0 {
		return fmt.Errorf("Integer arguments cannot be negative")
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		return err
	}

	var approvalP []byte
	if approvalProgram != "" {
		approvalP, err = ioutil.ReadFile("features/resources/" + approvalProgram)
		if err != nil {
			return err
		}
	}

	var clearP []byte
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

	fAssets, err := splitUint64(foreignAssets)
	if err != nil {
		return err
	}

	gSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	lSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}

	gh, err := base64.StdEncoding.DecodeString(genesisHash)
	if err != nil {
		return err
	}

	// this is only kept to keep compatability with old features
	// going forward, use txnSuggestedParams
	suggestedParams := types.SuggestedParams{
		Fee:             types.MicroAlgos(uint64(fee)),
		GenesisID:       "",
		GenesisHash:     gh,
		FirstRoundValid: types.Round(firstValid),
		LastRoundValid:  types.Round(lastValid),
		FlatFee:         true,
	}

	switch operation {
	case "create":
		tx, err = future.MakeApplicationCreateTxWithExtraPages(false, approvalP, clearP,
			gSchema, lSchema, args, accs, fApp, fAssets,
			suggestedParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{}, uint32(extraPages))
	case "update":
		tx, err = future.MakeApplicationUpdateTx(applicationId, args, accs, fApp, fAssets,
			approvalP, clearP,
			suggestedParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "call":
		tx, err = future.MakeApplicationCallTx(applicationId, args, accs,
			fApp, fAssets, types.NoOpOC, approvalP, clearP, gSchema, lSchema,
			suggestedParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "optin":
		tx, err = future.MakeApplicationOptInTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "clear":
		tx, err = future.MakeApplicationClearStateTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "closeout":
		tx, err = future.MakeApplicationCloseOutTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "delete":
		tx, err = future.MakeApplicationDeleteTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, senderAddr, nil, types.Digest{}, [32]byte{}, types.Address{})
	default:
		err = fmt.Errorf("Unknown opperation: %s", operation)
	}
	return err

}

func TransactionsUnitContext(s *godog.Suite) {
	// @unit.transactions
	s.Step(`^a signing account with address "([^"]*)" and mnemonic "([^"]*)"$`, aSigningAccountWithAddressAndMnemonic)
	s.Step(`^suggested transaction parameters fee (\d+), flat-fee "([^"]*)", first-valid (\d+), last-valid (\d+), genesis-hash "([^"]*)", genesis-id "([^"]*)"$`, suggestedTransactionParameters)
	s.Step(`^sign the transaction$`, signTheTransaction)
	s.Step(`^the base64 encoded signed transaction should equal "([^"]*)"$`, theBase64EncodedSignedTransactionShouldEqual)
	s.Step(`^the decoded transaction should equal the original$`, theDecodedTransactionShouldEqualTheOriginal)

	// @unit.transactions.keyreg
	s.Step(`^I build a keyreg transaction with sender "([^"]*)", nonparticipation "([^"]*)", vote first (\d+), vote last (\d+), key dilution (\d+), vote public key "([^"]*)", selection public key "([^"]*)", and state proof public key "([^"]*)"$`, buildKeyregTransaction)

	s.Step(`^I build an application transaction with operation "([^"]*)", application-id (\d+), sender "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", foreign-assets "([^"]*)", app-accounts "([^"]*)", fee (\d+), first-valid (\d+), last-valid (\d+), genesis-hash "([^"]*)", extra-pages (\d+)$`, buildLegacyAppCallTransaction)
}
